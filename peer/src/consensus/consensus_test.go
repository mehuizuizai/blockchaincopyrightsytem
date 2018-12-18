package consensus

//"tesing"

func TestDoVote(t *testing.T) {
	conSessionMap["123"]["192.168.13.81"] = []byte("1")
	conSessionMap["123"]["192.168.13.79"] = []byte("1")
	conSessionMap["123"]["192.168.13.65"] = []byte("1")

	isSucessful, mostVote, mostPeers := doVote("123")

	t.Println("isSucessful=", isSucessful)
	t.Println("mostVote=", mostVote)
	t.Println("mostPeers=", mostPeers)
}
