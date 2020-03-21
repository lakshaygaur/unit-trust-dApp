package main

type Fund struct {
	FundId		string		`json:"fundId"`
	Type		string		`json:"type"`
	Value		string		`json:"value"`
	ValidFrom	string		`json:"validFrom"`
	ValidTo		string 		`json:"validTo"`
	Owner		string		`json:"owner"`
	// TxnHistory  []TransactionHistory	`json:"txnHistory"`
}


type TransactionHistory struct {
	TxnId		string		`json:"txnId"`
	Status		string		`json:"status"`
	Timestamp	string		`json:"timestamp"`
}


type Account struct {
	AccountId	string		`json:"accountId"`
	Name		string		`json:"name"`
	Type		string		`json:"type"`
	Status		bool		`json:"status"`
}