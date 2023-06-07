package db

import "fmt"

type Stats struct {
	ID         string
	CountGames int
	CountWins  int
	TotalTime  int
}

func (m *Manager) CreateStatsTable() error {
	_, err := m.DB.Exec(`CREATE TABLE IF NOT EXISTS stats (
    	id text PRIMARY KEY,
    	count_games integer,
    	count_wins integer,
    	time integer
   	);`)
	if err != nil {
		return fmt.Errorf("failed to create table stats: %v", err)
	}
	return nil
}

func (m *Manager) GetStats(ID string) (*Stats, error) {
	rows, err := m.DB.Query("SELECT * FROM stats WHERE stats.id=?", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	stats := &Stats{}
	err = rows.Scan(&stats.ID, &stats.CountGames, &stats.CountWins, &stats.TotalTime)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %v", err)
	}
	return stats, nil
}

func (m *Manager) SelectStats() ([]*Stats, error) {
	rows, err := m.DB.Query("SELECT * FROM stats")
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}
	defer rows.Close()

	var listStats []*Stats
	for rows.Next() {
		stats := &Stats{}
		err = rows.Scan(&stats.ID, &stats.CountGames, &stats.CountWins, &stats.TotalTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}
		listStats = append(listStats, stats)
	}
	return listStats, nil
}

func (m *Manager) CreateStats(stats *Stats) error {
	_, err := m.DB.Exec(
		"INSERT INTO stats VALUES (?, ?, ?, ?)",
		stats.ID, stats.CountGames, stats.CountWins, stats.TotalTime,
	)
	if err != nil {
		return fmt.Errorf("failed to create stats: %v", err)
	}
	return nil
}

func (m *Manager) UpdateStats(stats *Stats) error {
	_, err := m.DB.Exec(
		"UPDATE stats SET VALUES (stats.id, stats.count_games, stats.count_wins, stats.time) = (?, ?, ?, ?) WHERE stats.id=?",
		stats.ID, stats.CountGames, stats.CountWins, stats.TotalTime,
	)
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}
	return nil
}

func (m *Manager) DeleteStats(stats *Stats) error {
	_, err := m.DB.Exec("DELETE FROM stats WHERE stats.id=?", stats.ID)
	if err != nil {
		return fmt.Errorf("failed to delete stats: %v", err)
	}
	return nil
}
