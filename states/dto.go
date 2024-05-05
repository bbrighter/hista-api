package states

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

func (states States) toResponse() StatesResponse {
	var resp []StateMetaResponse
	for _, s := range states {
		resp = append(resp,
			StateMetaResponse{
				ID:   s.ID,
				Date: s.Date,
			})
	}
	return StatesResponse{States: resp}
}

func (state State) toResponse() StateResponse {
	var conditionResp []ConditionResponse
	for _, con := range state.Conditions {
		conditionResp = append(conditionResp, con.toResponse())
	}
	var resp = StateResponse{
		ID:         state.ID,
		Date:       state.Date,
		Conditions: conditionResp,
	}
	return resp
}
