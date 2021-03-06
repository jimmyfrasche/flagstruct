//Package flagstruct allows flags to be defined by struct tags and automatically
//filled into a struct.
//
//For a flag to be registered the field must be exported, of a supported type,
//and have a flag struct tag.
//
//The package example covers most variations.
//
//Flag struct tag format
//
//The flag struct tag format is
//	flag:"name,default-value,description"
//
//If name is omitted, the field's name is used, lowercased, with _ replaced by -.
//
//If default-value is omitted, the appropriate zero value is used.
//
//If description is omitted, the empty string is used.
//
//If all are omitted the tag may be just:
//	flag:""
//
//If you need to use commas in the default-value, you can change the separator
//used by <rune>:<specification>.
//For example,
//	flag:"|:name|a,b,c|description"
//	flag:"::name:a,b,c:description"
//	flag:"⊕:name⊕a,b,c⊕description"
//
//Supported Types
//
//Supported types are:
//	string
//	bool
//	int
//	int64
//	uint
//	uint64
//	float64
//
//Or any type whose underlying type is one of the above.
//A type whose underlying type is one of the above is treated as that type,
//with the sole exception of a time.Duration which is handled as a time.Duration.
//
//Types whose field is exported are scanned for flags as well.
package flagstruct
