// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD4176298DecodeGithubComGoParkMailRu20211YSNPInternalAppModels(in *jlexer.Lexer, out *Search) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Category":
			out.Category = string(in.String())
		case "Date":
			out.Date = string(in.String())
		case "FromAmount":
			out.FromAmount = int(in.Int())
		case "ToAmount":
			out.ToAmount = int(in.Int())
		case "Radius":
			out.Radius = uint64(in.Uint64())
		case "Latitude":
			out.Latitude = float64(in.Float64())
		case "Longitude":
			out.Longitude = float64(in.Float64())
		case "Search":
			out.Search = string(in.String())
		case "Sorting":
			out.Sorting = string(in.String())
		case "From":
			out.From = uint64(in.Uint64())
		case "Count":
			out.Count = uint64(in.Uint64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD4176298EncodeGithubComGoParkMailRu20211YSNPInternalAppModels(out *jwriter.Writer, in Search) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Category\":"
		out.RawString(prefix[1:])
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"Date\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"FromAmount\":"
		out.RawString(prefix)
		out.Int(int(in.FromAmount))
	}
	{
		const prefix string = ",\"ToAmount\":"
		out.RawString(prefix)
		out.Int(int(in.ToAmount))
	}
	{
		const prefix string = ",\"Radius\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Radius))
	}
	{
		const prefix string = ",\"Latitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Latitude))
	}
	{
		const prefix string = ",\"Longitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Longitude))
	}
	{
		const prefix string = ",\"Search\":"
		out.RawString(prefix)
		out.String(string(in.Search))
	}
	{
		const prefix string = ",\"Sorting\":"
		out.RawString(prefix)
		out.String(string(in.Sorting))
	}
	{
		const prefix string = ",\"From\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.From))
	}
	{
		const prefix string = ",\"Count\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Search) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeGithubComGoParkMailRu20211YSNPInternalAppModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Search) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeGithubComGoParkMailRu20211YSNPInternalAppModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Search) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeGithubComGoParkMailRu20211YSNPInternalAppModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Search) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeGithubComGoParkMailRu20211YSNPInternalAppModels(l, v)
}