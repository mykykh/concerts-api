package domain

import (
    "time"
)

type concert struct {
    ID int64

    Name string
    Description string
    CreateTime time.Time
    UpdateTime time.Time
}
