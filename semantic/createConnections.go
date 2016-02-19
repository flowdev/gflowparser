package main


type SemanticConnections struct {
 createConns *CreateConnections
 verifyOutPorts *VerifyOutPortsUsedOnlyOnce
 handleChainBeg *HandleChainBeg
 handleChainMids *HandleChainMids
 handleChainEnd *HandleChainEnd
 begAddLastOp *AddLastOp
 midAddLastOp *AddLastOp
	InPort func(ParserData)
}

func NewSemanticConnections() *SemanticConnections {
    f := &SemanticConnections{}
 f.createConns = NewCreateConnections()
 f.verifyOutPorts = NewVerifyOutPortsUsedOnlyOnce()
 f.handleChainBeg = NewHandleChainBeg()
 f.handleChainMids = NewHandleChainMids()
 f.handleChainEnd = NewHandleChainEnd()
 f.begAddLastOp = NewAddLastOp()
 f.midAddLastOp = NewAddLastOp()

    f.createConns.SetOutPort(f.verifyOutPorts.InPort)
    f.createConns.SetChainOutPort(f.handleChainBeg.InPort)
    f.handleChainBeg.SetOutPort(f.handleChainMids.InPort)
    f.handleChainMids.SetOutPort(f.handleChainEnd.InPort)
    f.handleChainEnd.SetOutPort(f.createConns.ChainInPort)
    f.handleChainBeg.SetAddOpOutPort(f.begAddLastOp.InPort)
    f.begAddLastOp.SetOutPort(f.handleChainBeg.AddOpInPort)
    f.handleChainMids.SetAddOpOutPort(f.midAddLastOp.InPort)
    f.midAddLastOp.SetOutPort(f.handleChainMids.AddOpInPort)

	f.InPort = f.createConns.InPort

    return f
}
func (f *SemanticConnections) SetOutPort(port func()) {
	f.verifyOutPorts.SetOutPort(port)
}

