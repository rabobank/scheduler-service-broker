package model

import "time"

type GenericV3Response struct {
	Pagination struct {
		TotalResults int `json:"total_results"`
	} `json:"pagination"`
	Resources []struct {
		Guid string `json:"guid"`
	} `json:"resources"`
}

type TaskListResponse struct {
	Pagination struct {
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
		First        struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Previous interface{} `json:"previous"`
	} `json:"pagination"`
	Resources []struct {
		GUID       string    `json:"guid"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		SequenceID int       `json:"sequence_id"`
		Name       string    `json:"name"`
		State      string    `json:"state"`
		MemoryInMb int       `json:"memory_in_mb"`
		DiskInMb   int       `json:"disk_in_mb"`
		Result     struct {
			FailureReason string `json:"failure_reason"`
		} `json:"result"`
		DropletGUID   string `json:"droplet_guid"`
		Relationships struct {
			App struct {
				Data struct {
					GUID string `json:"guid"`
				} `json:"data"`
			} `json:"app"`
		} `json:"relationships"`
		Metadata struct {
			Labels struct {
			} `json:"labels"`
			Annotations struct {
			} `json:"annotations"`
		} `json:"metadata"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			App struct {
				Href string `json:"href"`
			} `json:"app"`
			Cancel struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"cancel"`
			Droplet struct {
				Href string `json:"href"`
			} `json:"droplet"`
		} `json:"links"`
	} `json:"resources"`
}

type TaskCreateResponse struct {
	GUID       string `json:"guid"`
	SequenceID int    `json:"sequence_id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	State      string `json:"state"`
	MemoryInMb int    `json:"memory_in_mb"`
	DiskInMb   int    `json:"disk_in_mb"`
	Result     struct {
		FailureReason interface{} `json:"failure_reason"`
	} `json:"result"`
	DropletGUID string `json:"droplet_guid"`
	Metadata    struct {
		Labels struct {
		} `json:"labels"`
		Annotations struct {
		} `json:"annotations"`
	} `json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Relationships struct {
		App struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
	} `json:"relationships"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		App struct {
			Href string `json:"href"`
		} `json:"app"`
		Cancel struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"cancel"`
		Droplet struct {
			Href string `json:"href"`
		} `json:"droplet"`
	} `json:"links"`
}
