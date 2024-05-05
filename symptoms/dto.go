package symptoms

func (symptom Symptom) toResponse() SymptomResponse {
	return SymptomResponse{
		ID:   symptom.ID,
		Name: symptom.Name,
	}
}

func (cats SymptomCategories) toResponse() SymptomCategoriesResponse {
	var resp []SymptomCategoryResponse
	for _, cat := range cats {
		var symptomsResp []SymptomResponse
		for _, sym := range cat.Symptoms {
			symptomsResp = append(symptomsResp, sym.toResponse())
		}
		var c = SymptomCategoryResponse{
			ID:       cat.ID,
			Name:     cat.Name,
			Symptoms: symptomsResp,
		}
		resp = append(resp, c)
	}
	return SymptomCategoriesResponse{Categories: resp}
}

func (condition Condition) toResponse() ConditionResponse {
	return ConditionResponse{
		ID:       condition.ID,
		Symptom:  condition.Symptom.toResponse(),
		Severity: condition.Severity,
	}
}

func (events ConditionEvents) toResponse() ConditionEventsResponse {
	var resp []ConditionEventMetaResponse
	for _, s := range events {
		resp = append(resp,
			ConditionEventMetaResponse{
				ID:   s.ID,
				Date: s.Date,
			})
	}
	return ConditionEventsResponse{ConditionEvents: resp}
}

func (event ConditionEvent) toResponse() ConditionEventResponse {
	var conditionResp []ConditionResponse
	for _, con := range event.Conditions {
		conditionResp = append(conditionResp, con.toResponse())
	}
	var resp = ConditionEventResponse{
		ID:         event.ID,
		Date:       event.Date,
		Conditions: conditionResp,
	}
	return resp
}
