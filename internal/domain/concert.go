package domain

import (
    "time"
)

type Concert struct {
    ID int64 `json:"id"`

    Title string `json:"title"`
    Description string `json:"description"`
    Location string `json:"location"`
    CreateDate time.Time `json:"create-date"`
    UpdateDate time.Time `json:"update-date"`
}
