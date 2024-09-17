package model

import (
	"github.com/abdelrahman146/digital-wallet/pkg/types"
)

type Auditable struct {
	actor     string
	actorId   string
	remarks   *string
	oldRecord interface{}
}

func (a *Auditable) SetActor(actor string, actorId string) {
	a.actor = actor
	a.actorId = actorId
}

func (a *Auditable) GetActor() (actor string, actorId string) {
	return a.actor, a.actorId
}

func (a *Auditable) SetRemarks(remarks string) {
	a.remarks = &remarks
}

func (a *Auditable) GetRemarks() *string {
	return a.remarks
}

func (a *Auditable) SetOldRecord(oldRecord interface{}) {
	a.oldRecord = oldRecord
}

func (a *Auditable) GetOldRecord() interface{} {
	return a.oldRecord
}

func (a *Auditable) CreateAudit(table, operation, recordId string, newRecord interface{}) (*Audit, error) {
	var oldRecordJSON types.JSONB
	var newRecordJSON types.JSONB
	audit := &Audit{
		Actor:     a.actor,
		ActorID:   a.actorId,
		Table:     table,
		Operation: operation,
		RecordID:  recordId,
		Remarks:   a.remarks,
	}
	if a.oldRecord != nil {
		if err := types.StructToJSONB(a.oldRecord, &oldRecordJSON); err != nil {
			return nil, err
		}
		audit.OldRecord = &oldRecordJSON
	}
	if newRecord != nil {
		if err := types.StructToJSONB(newRecord, &newRecordJSON); err != nil {
			return nil, err
		}
		audit.NewRecord = &newRecordJSON
	}
	return audit, nil
}
