package inserter

type Requester interface {
}

type MakeReq func(q, content string, count int)
