package chip8

const (
	System                 uint16 = 0x0000
	Clear                         = 0x00e0
	Return                        = 0x00ee
	Jump                          = 0x1000
	Call                          = 0x2000
	SkipIfEqual                   = 0x3000
	SkipIfNotEqual                = 0x4000
	SkipIfEqualRegister           = 0x5000
	SetValue                      = 0x6000
	AddValue                      = 0x7000
	SetRegister                   = 0x8000
	Or                            = 0x8001
	And                           = 0x8002
	Xor                           = 0x8003
	AddRegister                   = 0x8004
	SubtractYFromX                = 0x8005
	ShiftRight                    = 0x8006
	SubtractXFromY                = 0x8007
	ShiftLeft                     = 0x800e
	SkipIfNotEqualRegister        = 0x9000
	SetIndex                      = 0xa000
	JumpRelative                  = 0xb000
	AndRandom                     = 0xc000
	Draw                          = 0xd000
	SkipIfKeyPressed              = 0xe09e
	SkipIfKeyNotPressed           = 0xe0a1
	StoreDelayTimer               = 0xf007
	AwaitKeyPress                 = 0xf00a
	SetDelayTimer                 = 0xf015
	SetSoundTimer                 = 0xf018
	AddIndex                      = 0xf01e
	SetIndexFontCharacter         = 0xf029
	StoreBCD                      = 0xf033
	WriteMemory                   = 0xf055
	ReadMemory                    = 0xf065
)
