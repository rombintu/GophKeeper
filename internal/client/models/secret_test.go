package models

import (
	"testing"

	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

func TestSecretText_GetTitle(t *testing.T) {
	st := SecretText{
		Secret: Secret{
			Title: "Test Title",
			Type:  kpb.Secret_TEXT,
		},
		Text: "This is a test text",
	}

	if st.GetTitle() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", st.GetTitle())
	}
}

func TestSecretText_GetType(t *testing.T) {
	st := SecretText{
		Secret: Secret{
			Title: "Test Title",
			Type:  kpb.Secret_TEXT,
		},
		Text: "This is a test text",
	}

	if st.GetType() != kpb.Secret_TEXT {
		t.Errorf("Expected type 'TEXT', got '%v'", st.GetType())
	}
}

func TestSecretText_Payload(t *testing.T) {
	st := SecretText{
		Secret: Secret{
			Title: "Test Title",
			Type:  kpb.Secret_TEXT,
		},
		Text: "This is a test text",
	}

	expectedPayload := []byte("This is a test text")
	if string(st.Payload()) != string(expectedPayload) {
		t.Errorf("Expected payload '%s', got '%s'", expectedPayload, st.Payload())
	}
}

func TestSecretCreds_GetTitle(t *testing.T) {
	sc := SecretCreds{
		Secret: Secret{
			Title: "Test Creds",
			Type:  kpb.Secret_CRED,
		},
		Creds: Creds{
			URL:      "https://example.com",
			Login:    "user",
			Password: "password",
		},
	}

	if sc.GetTitle() != "Test Creds" {
		t.Errorf("Expected title 'Test Creds', got '%s'", sc.GetTitle())
	}
}

func TestSecretCreds_GetType(t *testing.T) {
	sc := SecretCreds{
		Secret: Secret{
			Title: "Test Creds",
			Type:  kpb.Secret_CRED,
		},
		Creds: Creds{
			URL:      "https://example.com",
			Login:    "user",
			Password: "password",
		},
	}

	if sc.GetType() != kpb.Secret_CRED {
		t.Errorf("Expected type 'CREDS', got '%v'", sc.GetType())
	}
}

func TestSecretCreds_Payload(t *testing.T) {
	sc := SecretCreds{
		Secret: Secret{
			Title: "Test Creds",
			Type:  kpb.Secret_CRED,
		},
		Creds: Creds{
			URL:      "https://example.com",
			Login:    "user",
			Password: "password",
		},
	}

	expectedPayload := `{"URL":"https://example.com","Login":"user","Password":"password"}`
	if string(sc.Payload()) != expectedPayload {
		t.Errorf("Expected payload '%s', got '%s'", expectedPayload, sc.Payload())
	}
}

func TestSecretBinary_GetTitle(t *testing.T) {
	sb := SecretBinary{
		Secret: Secret{
			Title: "Test Binary",
			Type:  kpb.Secret_DATA,
		},
		BinaryData: []byte{0x01, 0x02, 0x03},
	}

	if sb.GetTitle() != "Test Binary" {
		t.Errorf("Expected title 'Test Binary', got '%s'", sb.GetTitle())
	}
}

func TestSecretBinary_GetType(t *testing.T) {
	sb := SecretBinary{
		Secret: Secret{
			Title: "Test Binary",
			Type:  kpb.Secret_DATA,
		},
		BinaryData: []byte{0x01, 0x02, 0x03},
	}

	if sb.GetType() != kpb.Secret_DATA {
		t.Errorf("Expected type 'BINARY', got '%v'", sb.GetType())
	}
}

func TestSecretBinary_Payload(t *testing.T) {
	sb := SecretBinary{
		Secret: Secret{
			Title: "Test Binary",
			Type:  kpb.Secret_DATA,
		},
		BinaryData: []byte{0x01, 0x02, 0x03},
	}

	expectedPayload := []byte{0x01, 0x02, 0x03}
	if string(sb.Payload()) != string(expectedPayload) {
		t.Errorf("Expected payload '%v', got '%v'", expectedPayload, sb.Payload())
	}
}

func TestSecretCard_GetTitle(t *testing.T) {
	sc := SecretCard{
		Secret: Secret{
			Title: "Test Card",
			Type:  kpb.Secret_CARD,
		},
		Card: Card{
			Owner:      "John Doe",
			ExpireDate: "12/25",
			Number:     "1234 5678 9012 3456",
			Code:       "123",
		},
	}

	if sc.GetTitle() != "Test Card" {
		t.Errorf("Expected title 'Test Card', got '%s'", sc.GetTitle())
	}
}

func TestSecretCard_GetType(t *testing.T) {
	sc := SecretCard{
		Secret: Secret{
			Title: "Test Card",
			Type:  kpb.Secret_CARD,
		},
		Card: Card{
			Owner:      "John Doe",
			ExpireDate: "12/25",
			Number:     "1234 5678 9012 3456",
			Code:       "123",
		},
	}

	if sc.GetType() != kpb.Secret_CARD {
		t.Errorf("Expected type 'CARD', got '%v'", sc.GetType())
	}
}

func TestSecretCard_Payload(t *testing.T) {
	sc := SecretCard{
		Secret: Secret{
			Title: "Test Card",
			Type:  kpb.Secret_CARD,
		},
		Card: Card{
			Owner:      "John Doe",
			ExpireDate: "12/25",
			Number:     "1234 5678 9012 3456",
			Code:       "123",
		},
	}

	expectedPayload := `{"Owner":"John Doe","ExpireDate":"12/25","Number":"1234 5678 9012 3456","Code":"123"}`
	if string(sc.Payload()) != expectedPayload {
		t.Errorf("Expected payload '%s', got '%s'", expectedPayload, sc.Payload())
	}
}
