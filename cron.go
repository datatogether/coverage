package main

import (
	"database/sql"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/coverage"
	"time"
)

var lastUpdate time.Time

func cron() (stop chan (bool)) {
	stop = make(chan (bool), 0)

	go func() {
		for {
			select {
			case <-time.Tick(time.Hour * 6):
				if time.Since(lastUpdate) >= time.Hour {
					update(appDB)
				}
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func update(db *sql.DB) error {
	if cfg.RunCron {
		lastUpdate = time.Now()
		log.Info("updating:", lastUpdate)

		if err := calcSourceCoverage(db); err != nil {
			return err
		}

		if err := calcPrimerSourceCoverage(db); err != nil {
			return err
		}
	}

	return nil
}

func calcSourceCoverage(db *sql.DB) error {
	cvg := coverage.NewCoverageGenerator(nil, nil)
	pageSize := 100

	count, err := archive.CountSources(appDB)
	if err != nil {
		return err
	}

	numPages := count / pageSize
	for page := 0; page <= numPages; page++ {
		sources, err := archive.ListSources(db, pageSize, pageSize*page)
		if err != nil {
			return err
		}

		for _, s := range sources {
			summary, err := cvg.Summary()
			if err != nil {
				return err
			}

			if s.Stats == nil {
				s.Stats = &archive.SourceStats{}
			}

			if s.Stats.ArchivedUrlCount != summary.Archived {
				s.Stats.ArchivedUrlCount = summary.Archived
				s.Stats.UrlCount = summary.Descendants
				log.Infof("updating source: %s - %s: %f%%", s.Id, s.Title, float32(summary.Archived)/float32(summary.Descendants)*100)
				if err := s.Save(db); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func calcPrimerSourceCoverage(db *sql.DB) error {
	pageSize := 100

	count, err := archive.CountPrimers(appDB)
	if err != nil {
		return err
	}

	numPages := int(count) / pageSize
	for page := 0; page <= numPages; page++ {
		primers, err := archive.ListPrimers(db, pageSize, pageSize*page)
		if err != nil {
			return err
		}

		for _, p := range primers {
			if err := p.ReadSources(appDB); err != nil {
				return err
			}

			urlCount := 0
			archivedCount := 0

			for _, s := range p.Sources {
				if s.Stats != nil {
					urlCount = urlCount + s.Stats.UrlCount
					archivedCount = archivedCount + s.Stats.ArchivedUrlCount
				}
			}

			if p.Stats == nil {
				p.Stats = &archive.PrimerStats{}
			}

			if p.Stats.SourcesUrlCount != urlCount || p.Stats.SourcesArchivedUrlCount != archivedCount {
				p.Stats.SourcesUrlCount = urlCount
				p.Stats.SourcesArchivedUrlCount = archivedCount
				log.Infof("updating primer sources: %s - %s: %f%%", p.Id, p.ShortTitle, float32(p.Stats.SourcesArchivedUrlCount)/float32(p.Stats.SourcesUrlCount)*100)
				if err := p.Save(appDB); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// TODO - finish
func calcPrimerCoverage(db *sql.DB, primers []*archive.Primer) error {
	for _, primer := range primers {
		if err := primer.ReadSubPrimers(db); err != nil {
			return err
		}

		if len(primer.SubPrimers) == 0 && primer.Stats != nil {
			primer.Stats.UrlCount = primer.Stats.SourcesUrlCount
			primer.Stats.ArchivedUrlCount = primer.Stats.SourcesArchivedUrlCount
		}

		if err := calcPrimerCoverage(db, primer.SubPrimers); err != nil {
			return err
		}
	}
	return nil
}
