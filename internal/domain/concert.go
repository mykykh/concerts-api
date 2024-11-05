package domain

import (
    "time"
)

type Concert struct {
    ID int64

    Title string
    Description string
    Location string
    CreateDate time.Time
    UpdateDate time.Time
}
