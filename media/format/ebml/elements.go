package ebml

const (
	ElementUnknown            EBMLID = 0
	ElementVoid                      = 0xEC
	ElementCRC32                     = 0xBF
	ElementEBML                      = 0x1A45DFA3
	ElementEBMLVersion               = 0x4286
	ElementEBMLReadVersion           = 0x42F7
	ElementEBMLMaxIDLength           = 0x42F2
	ElementEBMLMaxSizeLength         = 0x42F3
	ElementDocType                   = 0x4282
	ElementDocTypeVersion            = 0x4287
	ElementDocTypeReadVersion        = 0x4285

	ElementSegment      = 0x18538067
	ElementSeekHead     = 0x114D9B74
	ElementSeek         = 0x4DBB
	ElementSeekID       = 0x53AB
	ElementSeekPosition = 0x53AC

	ElementInfo                       = 0x1549A966
	ElementSegmentUID                 = 0x73A4
	ElementSegmentFilename            = 0x7384
	ElementPrevUID                    = 0x3CB923
	ElementPrevFilename               = 0x3C83AB
	ElementNextUID                    = 0x3EB923
	ElementNextFilename               = 0x3E83BB
	ElementSegmentFamily              = 0x4444
	ElementChapterTranslate           = 0x6924
	ElementChapterTranslateEditionUID = 0x69FC
	ElementChapterTranslateCodec      = 0x69BF
	ElementChapterTranslateID         = 0x69A5
	ElementTimecodeScale              = 0x2AD7B1
	ElementDuration                   = 0x4489
	ElementDateUTC                    = 0x4461
	ElementTitle                      = 0x7BA9
	ElementMuxingApp                  = 0x4D80
	ElementWritingApp                 = 0x5741

	ElementCluster           = 0x1F43B675
	ElementTimecode          = 0xE7
	ElementSilentTracks      = 0x5854
	ElementSilentTrackNumber = 0x58D7
	ElementPosition          = 0xA7
	ElementPrevSize          = 0xAB
	ElementSimpleBlock       = 0xA3
	ElementBlockGroup        = 0xA0
	ElementBlock             = 0xA1

	ElementTracks = 0x1654AE6B
	// Tracks omitted... too many values

	ElementCues                = 0x1C53BB6B
	ElementCuePoint            = 0xBB
	ElementCueTime             = 0xB3
	ElementCueTrackPositions   = 0xB7
	ElementCueTrack            = 0xF7
	ElementCueClusterPosition  = 0xF1
	ElementCueRelativePosition = 0xF0
	ElementCueDuration         = 0xB2
	ElementCueBlockNumber      = 0x5378
	ElementCueCodecState       = 0xEA
	ElementCueReference        = 0xDB
	ElementCueRefTime          = 0x96
)

var Elements = map[EBMLID]ElementMeta{
	ElementUnknown: ElementMeta{
		Name: "Unknown", // not part of spec
		Type: typeGuess,
	},

	ElementVoid: ElementMeta{
		Name: "Void",
		Type: typeBinary,
	},
	ElementCRC32: ElementMeta{
		Name: "CRC-32",
		Type: typeBinary,
	},

	ElementEBML: ElementMeta{
		Name: "EBML",
		Type: typeMaster,
	},
	ElementEBMLVersion: ElementMeta{
		Name: "EBMLVersion",
		Type: typeUnsignedInt,
	},
	ElementEBMLReadVersion: ElementMeta{
		Name: "EBMLReadVersion",
		Type: typeUnsignedInt,
	},
	ElementEBMLMaxIDLength: ElementMeta{
		Name: "EBMLMaxIDLength",
		Type: typeUnsignedInt,
	},
	ElementEBMLMaxSizeLength: ElementMeta{
		Name: "EBMLMaxSizeLength",
		Type: typeUnsignedInt,
	},
	ElementDocType: ElementMeta{
		Name: "DocType",
		Type: typeString,
	},
	ElementDocTypeVersion: ElementMeta{
		Name: "DocTypeVersion",
		Type: typeUnsignedInt,
	},
	ElementDocTypeReadVersion: ElementMeta{
		Name: "DocTypeReadVersion",
		Type: typeUnsignedInt,
	},

	ElementSegment: ElementMeta{
		Name: "Segment",
		Type: typeMaster,
	},

	ElementSeekHead: ElementMeta{
		Name: "SeekHead",
		Type: typeMaster,
	},
	ElementSeek: ElementMeta{
		Name: "Seek",
		Type: typeMaster,
	},
	ElementSeekID: ElementMeta{
		Name: "SeekID",
		Type: typeBinary,
	},
	ElementSeekPosition: ElementMeta{
		Name: "SeekPosition",
		Type: typeUnsignedInt,
	},

	ElementInfo: ElementMeta{
		Name: "Info",
		Type: typeMaster,
	},
	ElementSegmentUID: ElementMeta{
		Name: "SegmentUID",
		Type: typeBinary,
	},
	ElementSegmentFilename: ElementMeta{
		Name: "SegmentFilename",
		Type: typeUtf8String,
	},
	ElementPrevUID: ElementMeta{
		Name: "PrevUID",
		Type: typeBinary,
	},
	ElementPrevFilename: ElementMeta{
		Name: "PrevFilename",
		Type: typeUtf8String,
	},
	ElementNextUID: ElementMeta{
		Name: "NextUID",
		Type: typeBinary,
	},
	ElementNextFilename: ElementMeta{
		Name: "NextFilename",
		Type: typeUtf8String,
	},
	ElementSegmentFamily: ElementMeta{
		Name: "SegmentFamily",
		Type: typeBinary,
	},
	ElementChapterTranslate: ElementMeta{
		Name: "ChapterTranslate",
		Type: typeMaster,
	},
	ElementChapterTranslateEditionUID: ElementMeta{
		Name: "ChapterTranslateEditionUID",
		Type: typeUnsignedInt,
	},
	ElementChapterTranslateCodec: ElementMeta{
		Name: "ChapterTranslateCodec",
		Type: typeUnsignedInt,
	},
	ElementChapterTranslateID: ElementMeta{
		Name: "ChapterTranslateID",
		Type: typeBinary,
	},
	ElementTimecodeScale: ElementMeta{
		Name:         "TimecodeScale",
		Type:         typeUnsignedInt,
		DefaultValue: uint64(1000000),
	},
	ElementDuration: ElementMeta{
		Name: "Duration",
		Type: typeFloat,
	},
	ElementDateUTC: ElementMeta{
		Name: "DateUTC",
		Type: typeDate,
	},
	ElementTitle: ElementMeta{
		Name: "Title",
		Type: typeUtf8String,
	},
	ElementMuxingApp: ElementMeta{
		Name: "MuxingApp",
		Type: typeUtf8String,
	},
	ElementWritingApp: ElementMeta{
		Name: "WritingApp",
		Type: typeUtf8String,
	},

	ElementCluster: ElementMeta{
		Name: "Cluster",
		Type: typeMaster,
	},
	ElementTimecode: ElementMeta{
		Name: "Timecode",
		Type: typeUnsignedInt,
	},
	ElementSilentTracks: ElementMeta{
		Name: "SilentTracks",
		Type: typeMaster,
	},
	ElementSilentTrackNumber: ElementMeta{
		Name: "SilentTrackNumber",
		Type: typeUnsignedInt,
	},
	ElementPosition: ElementMeta{
		Name: "Position",
		Type: typeUnsignedInt,
	},
	ElementPrevSize: ElementMeta{
		Name: "PrevSize",
		Type: typeUnsignedInt,
	},
	ElementSimpleBlock: ElementMeta{
		Name: "SimpleBlock",
		Type: typeBinary,
	},
	ElementBlockGroup: ElementMeta{
		Name: "BlockGroup",
		Type: typeMaster,
	},
	ElementBlock: ElementMeta{
		Name: "Block",
		Type: typeBinary,
	},

	ElementTracks: ElementMeta{
		Name: "Tracks",
		Type: typeMaster,
	},

	ElementCues: ElementMeta{
		Name: "Cues",
		Type: typeMaster,
	},
	ElementCuePoint: ElementMeta{
		Name: "CuePoint",
		Type: typeMaster,
	},
	ElementCueTime: ElementMeta{
		Name: "CueTime",
		Type: typeUnsignedInt,
	},
	ElementCueTrackPositions: ElementMeta{
		Name: "CueTrackPositions",
		Type: typeMaster,
	},
	ElementCueTrack: ElementMeta{
		Name: "CueTrack",
		Type: typeUnsignedInt,
	},
	ElementCueClusterPosition: ElementMeta{
		Name: "CueClusterPosition",
		Type: typeUnsignedInt,
	},
	ElementCueRelativePosition: ElementMeta{
		Name: "CueRelativePosition",
		Type: typeUnsignedInt,
	},
	ElementCueDuration: ElementMeta{
		Name: "CueDuration",
		Type: typeUnsignedInt,
	},
	ElementCueBlockNumber: ElementMeta{
		Name: "CueBlockNumber",
		Type: typeUnsignedInt,
	},
	ElementCueCodecState: ElementMeta{
		Name: "CueCodecState",
		Type: typeUnsignedInt,
	},
	ElementCueReference: ElementMeta{
		Name: "CueReference",
		Type: typeMaster,
	},
	ElementCueRefTime: ElementMeta{
		Name: "CueRefTime",
		Type: typeUnsignedInt,
	},
}
