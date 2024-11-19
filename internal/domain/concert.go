package domain

import (
    "time"

    "github.com/gofrs/uuid"
)

type Concert struct {
    ID int64 `json:"id"`

    Title string `json:"title"`
    Description string `json:"description"`
    Location string `json:"location"`
    AuthorID uuid.UUID `json:"author-id"`
    CreateDate time.Time `json:"create-date"`
    UpdateDate time.Time `json:"update-date"`
}
