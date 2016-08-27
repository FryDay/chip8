package chip8

const (
	opSystem                 uint16 = 0x0000
	opClear                         = 0x00e0
	opReturn                        = 0x00ee
	opJump                          = 0x1000
	opCall                          = 0x2000
	opSkipIfEqual                   = 0x3000
	opSkipIfNotEqual                = 0x4000
	opSkipIfEqualRegister           = 0x5000
	opSetValue                      = 0x6000
	opAddValue                      = 0x7000
	opSetRegister                   = 0x8000
	opOr                            = 0x8001
	opAnd                           = 0x8002
	opXor                           = 0x8003
	opAddRegister                   = 0x8004
	opSubtractYFromX                = 0x8005
	opShiftRight                    = 0x8006
	opSubtractXFromY                = 0x8007
	opShiftLeft                     = 0x800e
	opSkipIfNotEqualRegister        = 0x9000
	opSetIndex                      = 0xa000
	opJumpRelative                  = 0xb000
	opAndRandom                     = 0xc000
	opDraw                          = 0xd000
	opSkipIfKeyPressed              = 0xe09e
	opSkipIfKeyNotPressed           = 0xe0a1
	opStoreDelayTimer               = 0xf007
	opAwaitKeyPress                 = 0xf00a
	opSetDelayTimer                 = 0xf015
	opSetSoundTimer                 = 0xf018
	opAddIndex                      = 0xf01e
	opSetIndexFontCharacter         = 0xf029
	opStoreBCD                      = 0xf033
	opWriteMemory                   = 0xf055
	opReadMemory                    = 0xf065
)
