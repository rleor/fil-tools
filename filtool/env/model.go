package env

import (
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/specs-storage/storage"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
	"io"
)

// ---- log
type Log struct {
	Timestamp uint64
	Trace     string // for errors

	Message string

	// additional data (Event info)
	Kind string
}

func (t *Log) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{164}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Timestamp (uint64) (uint64)
	if len("Timestamp") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Timestamp\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Timestamp"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Timestamp")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Timestamp)); err != nil {
		return err
	}

	// t.Trace (string) (string)
	if len("Trace") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Trace\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Trace"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Trace")); err != nil {
		return err
	}

	if len(t.Trace) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Trace was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Trace))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Trace)); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len("Message") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Message\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Message"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Message")); err != nil {
		return err
	}

	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Message))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Message)); err != nil {
		return err
	}

	// t.Kind (string) (string)
	if len("Kind") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Kind\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Kind"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Kind")); err != nil {
		return err
	}

	if len(t.Kind) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Kind was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Kind))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Kind)); err != nil {
		return err
	}
	return nil
}

func (t *Log) UnmarshalCBOR(r io.Reader) error {
	*t = Log{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("Log: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Timestamp (uint64) (uint64)
		case "Timestamp":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.Timestamp = uint64(extra)

			}
			// t.Trace (string) (string)
		case "Trace":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Trace = string(sval)
			}
			// t.Message (string) (string)
		case "Message":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Message = string(sval)
			}
			// t.Kind (string) (string)
		case "Kind":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Kind = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}

// ---- SectorPreCommitInfo
var lengthBufSectorPreCommitInfo = []byte{138}

type SectorPreCommitInfo struct {
	SealProof       abi.RegisteredSealProof
	SectorNumber    abi.SectorNumber
	SealedCID       cid.Cid `checked:"true"` // CommR
	SealRandEpoch   abi.ChainEpoch
	DealIDs         []abi.DealID
	Expiration      abi.ChainEpoch
	ReplaceCapacity bool // Whether to replace a "committed capacity" no-deal sector (requires non-empty DealIDs)
	// The committed capacity sector to replace, and it's deadline/partition location
	ReplaceSectorDeadline  uint64
	ReplaceSectorPartition uint64
	ReplaceSectorNumber    abi.SectorNumber
}

func (t *SectorPreCommitInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSectorPreCommitInfo); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.SealProof (abi.RegisteredSealProof) (int64)
	if t.SealProof >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SealProof)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.SealProof-1)); err != nil {
			return err
		}
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SealedCID (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.SealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.SealedCID: %w", err)
	}

	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	if t.SealRandEpoch >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SealRandEpoch)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.SealRandEpoch-1)); err != nil {
			return err
		}
	}

	// t.DealIDs ([]abi.DealID) (slice)
	if len(t.DealIDs) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.DealIDs was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.DealIDs))); err != nil {
		return err
	}
	for _, v := range t.DealIDs {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	// t.ReplaceCapacity (bool) (bool)
	if err := cbg.WriteBool(w, t.ReplaceCapacity); err != nil {
		return err
	}

	// t.ReplaceSectorDeadline (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.ReplaceSectorDeadline)); err != nil {
		return err
	}

	// t.ReplaceSectorPartition (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.ReplaceSectorPartition)); err != nil {
		return err
	}

	// t.ReplaceSectorNumber (abi.SectorNumber) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.ReplaceSectorNumber)); err != nil {
		return err
	}

	return nil
}

func (t *SectorPreCommitInfo) UnmarshalCBOR(r io.Reader) error {
	*t = SectorPreCommitInfo{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 10 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealProof = abi.RegisteredSealProof(extraI)
	}
	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.SealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SealedCID: %w", err)
		}

		t.SealedCID = c

	}
	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealRandEpoch = abi.ChainEpoch(extraI)
	}
	// t.DealIDs ([]abi.DealID) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.DealIDs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.DealIDs = make([]abi.DealID, extra)
	}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.DealIDs slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.DealIDs was not a uint, instead got %d", maj)
		}

		t.DealIDs[i] = abi.DealID(val)
	}

	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	// t.ReplaceCapacity (bool) (bool)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ReplaceCapacity = false
	case 21:
		t.ReplaceCapacity = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.ReplaceSectorDeadline (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorDeadline = uint64(extra)

	}
	// t.ReplaceSectorPartition (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorPartition = uint64(extra)

	}
	// t.ReplaceSectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorNumber = abi.SectorNumber(extra)

	}
	return nil
}

// ---- piece
type Piece struct {
	Piece    abi.PieceInfo
	DealInfo *api.PieceDealInfo // nil for pieces which do not appear in deals (e.g. filler pieces)
}

func (t *Piece) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{162}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Piece (abi.PieceInfo) (struct)
	if len("Piece") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Piece\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Piece"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Piece")); err != nil {
		return err
	}

	if err := t.Piece.MarshalCBOR(w); err != nil {
		return err
	}

	// t.DealInfo (api.PieceDealInfo) (struct)
	if len("DealInfo") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealInfo\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealInfo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealInfo")); err != nil {
		return err
	}

	if err := t.DealInfo.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *Piece) UnmarshalCBOR(r io.Reader) error {
	*t = Piece{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("Piece: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Piece (abi.PieceInfo) (struct)
		case "Piece":

			{

				if err := t.Piece.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.Piece: %w", err)
				}

			}
			// t.DealInfo (api.PieceDealInfo) (struct)
		case "DealInfo":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}
					t.DealInfo = new(api.PieceDealInfo)
					if err := t.DealInfo.UnmarshalCBOR(br); err != nil {
						return xerrors.Errorf("unmarshaling t.DealInfo pointer: %w", err)
					}
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}

// ---- sector info
type SectorInfo struct {
	State        SectorState
	SectorNumber abi.SectorNumber

	SectorType abi.RegisteredSealProof

	// Packing
	CreationTime int64 // unix seconds
	Pieces       []Piece

	// PreCommit1
	TicketValue   abi.SealRandomness
	TicketEpoch   abi.ChainEpoch
	PreCommit1Out storage.PreCommit1Out

	// PreCommit2
	CommD *cid.Cid
	CommR *cid.Cid // SectorKey
	Proof []byte

	PreCommitInfo    *SectorPreCommitInfo
	PreCommitDeposit big.Int
	PreCommitMessage *cid.Cid
	PreCommitTipSet  TipSetToken

	PreCommit2Fails uint64

	// WaitSeed
	SeedValue abi.InteractiveSealRandomness
	SeedEpoch abi.ChainEpoch

	// Committing
	CommitMessage *cid.Cid
	InvalidProofs uint64 // failed proof computations (doesn't validate with proof inputs; can't compute)

	// finalized times. after FinalizeSector event, will reset to zero.
	FinalizedTimes uint64

	// CCUpdate
	CCUpdate             bool
	CCPieces             []Piece
	UpdateSealed         *cid.Cid
	UpdateUnsealed       *cid.Cid
	ReplicaUpdateProof   storage.ReplicaUpdateProof
	ReplicaUpdateMessage *cid.Cid

	// Faults
	FaultReportMsg *cid.Cid

	// Recovery
	Return ReturnState

	// Termination
	TerminateMessage *cid.Cid
	TerminatedAt     abi.ChainEpoch

	// Debug
	LastErr string

	Log []Log

	// in recovering period.
	Recovering bool
}

func (t *SectorInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{184, 34}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.State (sealing.SectorState) (string)
	if len("State") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"State\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("State"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("State")); err != nil {
		return err
	}

	if len(t.State) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.State was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.State))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.State)); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)
	if len("SectorNumber") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"SectorNumber\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("SectorNumber"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("SectorNumber")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SectorType (abi.RegisteredSealProof) (int64)
	if len("SectorType") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"SectorType\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("SectorType"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("SectorType")); err != nil {
		return err
	}

	if t.SectorType >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SectorType)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.SectorType-1)); err != nil {
			return err
		}
	}

	// t.CreationTime (int64) (int64)
	if len("CreationTime") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CreationTime\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CreationTime"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CreationTime")); err != nil {
		return err
	}

	if t.CreationTime >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.CreationTime)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.CreationTime-1)); err != nil {
			return err
		}
	}

	// t.Pieces ([]sealing.Piece) (slice)
	if len("Pieces") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Pieces\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Pieces"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Pieces")); err != nil {
		return err
	}

	if len(t.Pieces) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Pieces was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Pieces))); err != nil {
		return err
	}
	for _, v := range t.Pieces {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}

	// t.TicketValue (abi.SealRandomness) (slice)
	if len("TicketValue") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"TicketValue\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("TicketValue"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("TicketValue")); err != nil {
		return err
	}

	if len(t.TicketValue) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.TicketValue was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.TicketValue))); err != nil {
		return err
	}

	if _, err := w.Write(t.TicketValue[:]); err != nil {
		return err
	}

	// t.TicketEpoch (abi.ChainEpoch) (int64)
	if len("TicketEpoch") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"TicketEpoch\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("TicketEpoch"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("TicketEpoch")); err != nil {
		return err
	}

	if t.TicketEpoch >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.TicketEpoch)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.TicketEpoch-1)); err != nil {
			return err
		}
	}

	// t.PreCommit1Out (storage.PreCommit1Out) (slice)
	if len("PreCommit1Out") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommit1Out\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommit1Out"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommit1Out")); err != nil {
		return err
	}

	if len(t.PreCommit1Out) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.PreCommit1Out was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.PreCommit1Out))); err != nil {
		return err
	}

	if _, err := w.Write(t.PreCommit1Out[:]); err != nil {
		return err
	}

	// t.CommD (cid.Cid) (struct)
	if len("CommD") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CommD\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CommD"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CommD")); err != nil {
		return err
	}

	if t.CommD == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.CommD); err != nil {
			return xerrors.Errorf("failed to write cid field t.CommD: %w", err)
		}
	}

	// t.CommR (cid.Cid) (struct)
	if len("CommR") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CommR\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CommR"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CommR")); err != nil {
		return err
	}

	if t.CommR == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.CommR); err != nil {
			return xerrors.Errorf("failed to write cid field t.CommR: %w", err)
		}
	}

	// t.Proof ([]uint8) (slice)
	if len("Proof") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Proof\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Proof"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Proof")); err != nil {
		return err
	}

	if len(t.Proof) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Proof was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Proof))); err != nil {
		return err
	}

	if _, err := w.Write(t.Proof[:]); err != nil {
		return err
	}

	// t.PreCommitInfo (miner.SectorPreCommitInfo) (struct)
	if len("PreCommitInfo") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommitInfo\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommitInfo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommitInfo")); err != nil {
		return err
	}

	if err := t.PreCommitInfo.MarshalCBOR(w); err != nil {
		return err
	}

	// t.PreCommitDeposit (big.Int) (struct)
	if len("PreCommitDeposit") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommitDeposit\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommitDeposit"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommitDeposit")); err != nil {
		return err
	}

	if err := t.PreCommitDeposit.MarshalCBOR(w); err != nil {
		return err
	}

	// t.PreCommitMessage (cid.Cid) (struct)
	if len("PreCommitMessage") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommitMessage\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommitMessage"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommitMessage")); err != nil {
		return err
	}

	if t.PreCommitMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.PreCommitMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.PreCommitMessage: %w", err)
		}
	}

	// t.PreCommitTipSet (sealing.TipSetToken) (slice)
	if len("PreCommitTipSet") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommitTipSet\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommitTipSet"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommitTipSet")); err != nil {
		return err
	}

	if len(t.PreCommitTipSet) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.PreCommitTipSet was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.PreCommitTipSet))); err != nil {
		return err
	}

	if _, err := w.Write(t.PreCommitTipSet[:]); err != nil {
		return err
	}

	// t.PreCommit2Fails (uint64) (uint64)
	if len("PreCommit2Fails") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PreCommit2Fails\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("PreCommit2Fails"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PreCommit2Fails")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.PreCommit2Fails)); err != nil {
		return err
	}

	// t.SeedValue (abi.InteractiveSealRandomness) (slice)
	if len("SeedValue") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"SeedValue\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("SeedValue"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("SeedValue")); err != nil {
		return err
	}

	if len(t.SeedValue) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.SeedValue was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.SeedValue))); err != nil {
		return err
	}

	if _, err := w.Write(t.SeedValue[:]); err != nil {
		return err
	}

	// t.SeedEpoch (abi.ChainEpoch) (int64)
	if len("SeedEpoch") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"SeedEpoch\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("SeedEpoch"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("SeedEpoch")); err != nil {
		return err
	}

	if t.SeedEpoch >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.SeedEpoch)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.SeedEpoch-1)); err != nil {
			return err
		}
	}

	// t.CommitMessage (cid.Cid) (struct)
	if len("CommitMessage") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CommitMessage\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CommitMessage"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CommitMessage")); err != nil {
		return err
	}

	if t.CommitMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.CommitMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.CommitMessage: %w", err)
		}
	}

	// t.InvalidProofs (uint64) (uint64)
	if len("InvalidProofs") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"InvalidProofs\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("InvalidProofs"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("InvalidProofs")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.InvalidProofs)); err != nil {
		return err
	}

	// t.CCUpdate (bool) (bool)
	if len("CCUpdate") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CCUpdate\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CCUpdate"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CCUpdate")); err != nil {
		return err
	}

	if err := cbg.WriteBool(w, t.CCUpdate); err != nil {
		return err
	}

	// t.CCPieces ([]sealing.Piece) (slice)
	if len("CCPieces") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CCPieces\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("CCPieces"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CCPieces")); err != nil {
		return err
	}

	if len(t.CCPieces) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.CCPieces was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.CCPieces))); err != nil {
		return err
	}
	for _, v := range t.CCPieces {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}

	// t.UpdateSealed (cid.Cid) (struct)
	if len("UpdateSealed") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"UpdateSealed\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("UpdateSealed"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("UpdateSealed")); err != nil {
		return err
	}

	if t.UpdateSealed == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.UpdateSealed); err != nil {
			return xerrors.Errorf("failed to write cid field t.UpdateSealed: %w", err)
		}
	}

	// t.UpdateUnsealed (cid.Cid) (struct)
	if len("UpdateUnsealed") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"UpdateUnsealed\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("UpdateUnsealed"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("UpdateUnsealed")); err != nil {
		return err
	}

	if t.UpdateUnsealed == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.UpdateUnsealed); err != nil {
			return xerrors.Errorf("failed to write cid field t.UpdateUnsealed: %w", err)
		}
	}

	// t.ReplicaUpdateProof (storage.ReplicaUpdateProof) (slice)
	if len("ReplicaUpdateProof") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ReplicaUpdateProof\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("ReplicaUpdateProof"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ReplicaUpdateProof")); err != nil {
		return err
	}

	if len(t.ReplicaUpdateProof) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.ReplicaUpdateProof was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.ReplicaUpdateProof))); err != nil {
		return err
	}

	if _, err := w.Write(t.ReplicaUpdateProof[:]); err != nil {
		return err
	}

	// t.ReplicaUpdateMessage (cid.Cid) (struct)
	if len("ReplicaUpdateMessage") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ReplicaUpdateMessage\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("ReplicaUpdateMessage"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ReplicaUpdateMessage")); err != nil {
		return err
	}

	if t.ReplicaUpdateMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.ReplicaUpdateMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.ReplicaUpdateMessage: %w", err)
		}
	}

	// t.FaultReportMsg (cid.Cid) (struct)
	if len("FaultReportMsg") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"FaultReportMsg\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("FaultReportMsg"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("FaultReportMsg")); err != nil {
		return err
	}

	if t.FaultReportMsg == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.FaultReportMsg); err != nil {
			return xerrors.Errorf("failed to write cid field t.FaultReportMsg: %w", err)
		}
	}

	// t.Return (sealing.ReturnState) (string)
	if len("Return") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Return\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Return"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Return")); err != nil {
		return err
	}

	if len(t.Return) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Return was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Return))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Return)); err != nil {
		return err
	}

	// t.TerminateMessage (cid.Cid) (struct)
	if len("TerminateMessage") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"TerminateMessage\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("TerminateMessage"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("TerminateMessage")); err != nil {
		return err
	}

	if t.TerminateMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCidBuf(scratch, w, *t.TerminateMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.TerminateMessage: %w", err)
		}
	}

	// t.TerminatedAt (abi.ChainEpoch) (int64)
	if len("TerminatedAt") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"TerminatedAt\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("TerminatedAt"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("TerminatedAt")); err != nil {
		return err
	}

	if t.TerminatedAt >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.TerminatedAt)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.TerminatedAt-1)); err != nil {
			return err
		}
	}

	// t.LastErr (string) (string)
	if len("LastErr") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"LastErr\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("LastErr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("LastErr")); err != nil {
		return err
	}

	if len(t.LastErr) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.LastErr was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.LastErr))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.LastErr)); err != nil {
		return err
	}

	// t.Log ([]sealing.Log) (slice)
	if len("Log") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Log\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Log"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Log")); err != nil {
		return err
	}

	if len(t.Log) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Log was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Log))); err != nil {
		return err
	}
	for _, v := range t.Log {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}

	// t.FinalizedTimes (uint64) (uint64)
	if len("FinalizedTimes") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"FinalizedTimes\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("FinalizedTimes"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("FinalizedTimes")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.FinalizedTimes)); err != nil {
		return err
	}

	// t.Recovering (bool) (bool)
	if len("Recovering") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Recovering\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Recovering"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Recovering")); err != nil {
		return err
	}

	if err := cbg.WriteBool(w, t.Recovering); err != nil {
		return err
	}

	return nil
}

func (t *SectorInfo) UnmarshalCBOR(r io.Reader) error {
	*t = SectorInfo{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("SectorInfo: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.State (sealing.SectorState) (string)
		case "State":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.State = SectorState(sval)
			}
			// t.SectorNumber (abi.SectorNumber) (uint64)
		case "SectorNumber":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.SectorNumber = abi.SectorNumber(extra)

			}
			// t.SectorType (abi.RegisteredSealProof) (int64)
		case "SectorType":
			{
				maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.SectorType = abi.RegisteredSealProof(extraI)
			}
			// t.CreationTime (int64) (int64)
		case "CreationTime":
			{
				maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.CreationTime = int64(extraI)
			}
			// t.Pieces ([]sealing.Piece) (slice)
		case "Pieces":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.Pieces: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.Pieces = make([]Piece, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v Piece
				if err := v.UnmarshalCBOR(br); err != nil {
					return err
				}

				t.Pieces[i] = v
			}

			// t.TicketValue (abi.SealRandomness) (slice)
		case "TicketValue":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.TicketValue: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.TicketValue = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.TicketValue[:]); err != nil {
				return err
			}
			// t.TicketEpoch (abi.ChainEpoch) (int64)
		case "TicketEpoch":
			{
				maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.TicketEpoch = abi.ChainEpoch(extraI)
			}
			// t.PreCommit1Out (storage.PreCommit1Out) (slice)
		case "PreCommit1Out":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.PreCommit1Out: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.PreCommit1Out = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.PreCommit1Out[:]); err != nil {
				return err
			}
			// t.CommD (cid.Cid) (struct)
		case "CommD":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.CommD: %w", err)
					}

					t.CommD = &c
				}

			}
			// t.CommR (cid.Cid) (struct)
		case "CommR":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.CommR: %w", err)
					}

					t.CommR = &c
				}

			}
			// t.Proof ([]uint8) (slice)
		case "Proof":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Proof: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Proof = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.Proof[:]); err != nil {
				return err
			}
			// t.PreCommitInfo (miner.SectorPreCommitInfo) (struct)
		case "PreCommitInfo":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}
					t.PreCommitInfo = new(SectorPreCommitInfo)
					if err := t.PreCommitInfo.UnmarshalCBOR(br); err != nil {
						return xerrors.Errorf("unmarshaling t.PreCommitInfo pointer: %w", err)
					}
				}

			}
			// t.PreCommitDeposit (big.Int) (struct)
		case "PreCommitDeposit":

			{

				if err := t.PreCommitDeposit.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.PreCommitDeposit: %w", err)
				}

			}
			// t.PreCommitMessage (cid.Cid) (struct)
		case "PreCommitMessage":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.PreCommitMessage: %w", err)
					}

					t.PreCommitMessage = &c
				}

			}
			// t.PreCommitTipSet (sealing.TipSetToken) (slice)
		case "PreCommitTipSet":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.PreCommitTipSet: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.PreCommitTipSet = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.PreCommitTipSet[:]); err != nil {
				return err
			}
			// t.PreCommit2Fails (uint64) (uint64)
		case "PreCommit2Fails":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.PreCommit2Fails = uint64(extra)

			}
			// t.SeedValue (abi.InteractiveSealRandomness) (slice)
		case "SeedValue":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.SeedValue: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.SeedValue = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.SeedValue[:]); err != nil {
				return err
			}
			// t.SeedEpoch (abi.ChainEpoch) (int64)
		case "SeedEpoch":
			{
				maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.SeedEpoch = abi.ChainEpoch(extraI)
			}
			// t.CommitMessage (cid.Cid) (struct)
		case "CommitMessage":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.CommitMessage: %w", err)
					}

					t.CommitMessage = &c
				}

			}
			// t.InvalidProofs (uint64) (uint64)
		case "InvalidProofs":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.InvalidProofs = uint64(extra)

			}
			// t.CCUpdate (bool) (bool)
		case "CCUpdate":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}
			if maj != cbg.MajOther {
				return fmt.Errorf("booleans must be major type 7")
			}
			switch extra {
			case 20:
				t.CCUpdate = false
			case 21:
				t.CCUpdate = true
			default:
				return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
			}
			// t.CCPieces ([]sealing.Piece) (slice)
		case "CCPieces":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.CCPieces: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.CCPieces = make([]Piece, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v Piece
				if err := v.UnmarshalCBOR(br); err != nil {
					return err
				}

				t.CCPieces[i] = v
			}

			// t.UpdateSealed (cid.Cid) (struct)
		case "UpdateSealed":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.UpdateSealed: %w", err)
					}

					t.UpdateSealed = &c
				}

			}
			// t.UpdateUnsealed (cid.Cid) (struct)
		case "UpdateUnsealed":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.UpdateUnsealed: %w", err)
					}

					t.UpdateUnsealed = &c
				}

			}
			// t.ReplicaUpdateProof (storage.ReplicaUpdateProof) (slice)
		case "ReplicaUpdateProof":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.ReplicaUpdateProof: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.ReplicaUpdateProof = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.ReplicaUpdateProof[:]); err != nil {
				return err
			}
			// t.ReplicaUpdateMessage (cid.Cid) (struct)
		case "ReplicaUpdateMessage":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.ReplicaUpdateMessage: %w", err)
					}

					t.ReplicaUpdateMessage = &c
				}

			}
			// t.FaultReportMsg (cid.Cid) (struct)
		case "FaultReportMsg":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.FaultReportMsg: %w", err)
					}

					t.FaultReportMsg = &c
				}

			}
			// t.Return (sealing.ReturnState) (string)
		case "Return":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Return = ReturnState(sval)
			}
			// t.TerminateMessage (cid.Cid) (struct)
		case "TerminateMessage":

			{

				b, err := br.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := br.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(br)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.TerminateMessage: %w", err)
					}

					t.TerminateMessage = &c
				}

			}
			// t.TerminatedAt (abi.ChainEpoch) (int64)
		case "TerminatedAt":
			{
				maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.TerminatedAt = abi.ChainEpoch(extraI)
			}
			// t.LastErr (string) (string)
		case "LastErr":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.LastErr = string(sval)
			}
			// t.Log ([]sealing.Log) (slice)
		case "Log":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.Log: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.Log = make([]Log, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v Log
				if err := v.UnmarshalCBOR(br); err != nil {
					return err
				}

				t.Log[i] = v
			}
			// t.FinalizedTimes (uint64) (uint64)
		case "FinalizedTimes":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.FinalizedTimes = uint64(extra)

			}

			// t.Recovering (bool) (bool)
		case "Recovering":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}
			if maj != cbg.MajOther {
				return fmt.Errorf("booleans must be major type 7")
			}
			switch extra {
			case 20:
				t.Recovering = false
			case 21:
				t.Recovering = true
			default:
				return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}

type ReturnState string
type SectorState string
type TipSetToken []byte
