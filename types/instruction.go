package types

import "github.com/yosefl20/solana-go-sdk/common"

type CompiledInstruction struct {
	ProgramIDIndex int
	Accounts       []int
	Data           []byte
}

type Instruction struct {
	ProgramID common.PublicKey
	Accounts  []AccountMeta
	Data      []byte
}

type AccountMeta struct {
	PubKey     common.PublicKey
	IsSigner   bool
	IsWritable bool
}
