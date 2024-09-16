package model

import "digital-wallet/pkg/types"

type Auditable struct {
	actor     string
	actorId   *string
	remarks   *string
	oldRecord interface{}
}

func (a *Auditable) SetActor(actor string, actorId *string) {
	a.actor = actor
	a.actorId = actorId
}

func (a *Auditable) GetActor() (actor string, actorId *string) {
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

func (a *Auditable) CreateAudit(operation string, recordId string, newRecord interface{}) (*Audit, error) {
	var oldRecordJSON *types.JSONB
	var newRecordJSON *types.JSONB
	if a.oldRecord != nil {
		if err := oldRecordJSON.StructToJSONB(a.oldRecord); err != nil {
			return nil, err
		}
	}
	if newRecord != nil {
		if err := newRecordJSON.StructToJSONB(newRecord); err != nil {
			return nil, err
		}
	}
	return &Audit{
		Actor:     a.actor,
		ActorID:   a.actorId,
		Operation: operation,
		RecordID:  recordId,
		Remarks:   a.remarks,
		OldRecord: oldRecordJSON,
		NewRecord: newRecordJSON,
	}, nil
}
