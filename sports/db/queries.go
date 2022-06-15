package db

const (
	sportsList = "list"
)

func getSportsQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id,
				name,
				advertised_start_time 
			FROM sports 
		`,
	}
}
