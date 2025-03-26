package model

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Meta struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

func (m *Meta) SetMeta(page, limit, totalData int) {
	m.Page = page
	m.Limit = limit
	m.TotalData = totalData

	if totalData == 0 {
		m.TotalPage = 1
		return
	}

	m.TotalPage = totalData / limit
	if totalData%limit > 0 {
		m.TotalPage++
	}

	if m.TotalPage == 0 {
		m.TotalPage = 1
	}
}
