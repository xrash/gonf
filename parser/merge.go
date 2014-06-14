package parser

type table map[string]*PairNode

func merge(p *PairNode) {
	mergePairNode(p, nil)
}

func mergePairNode(p *PairNode, t *table) {
	if p == nil {
		return
	}

	if t == nil {
		t = &table{}
	}

	table := *t

	if root, exists := table[p.Key.Value.Value]; exists {
		var values *ValuesNode

		if root.Value.Array != nil {
			values = root.Value.Array.Values
			for values.Values != nil {
				values = values.Values
			}
		} else {
			values = NewValuesNode(root.Value, nil)
			array := NewArrayNode(values)
			value := NewValueNode(nil, nil, array)
			root.Value = value
		}

		root.Pair = p.Pair
		values.Values = NewValuesNode(p.Value, nil)
	} else {
		table[p.Key.Value.Value] = p
	}

	mergePairNode(p.Pair, t)
	mergeValueNode(p.Value)
}

func mergeValueNode(v *ValueNode) {
	switch {
	case v.Array != nil:
		mergeValuesNode(v.Array.Values)
	case v.Table != nil:
		mergePairNode(v.Table.Pair, nil)
	}
}

func mergeValuesNode(v *ValuesNode) {
	values := v

	for values != nil {
		mergeValueNode(values.Value)
		values = values.Values
	}
}
