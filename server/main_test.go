package main

import "testing"

func TestSanitizeTitle_SpecialChars(t *testing.T) {
	// input with leading/trailing spaces, lots of disallowed chars, and repeats
	input := "  ///Crazy:::Title???***(((!!!***///   "
	// after sanitization we expect single underscores, no leading/trailing, no repeats
	want := "Crazy_Title"
	got := sanitizeTitle(input)

	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_EmptyString(t *testing.T) {
	input := ""
	want := "untitled"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_WhitespaceOnly(t *testing.T) {
	input := "    "
	want := "untitled"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_AllInvalidChars(t *testing.T) {
	input := "////....::::????****"
	want := "untitled"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_LeadingTrailingInvalid(t *testing.T) {
	input := "  hello world!  "
	want := "hello_world"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_MultipleConsecutiveInvalid(t *testing.T) {
	input := "foo---bar___baz"
	want := "foo_bar_baz"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_ValidTitle(t *testing.T) {
	input := "My Report 2025"
	want := "My_Report_2025"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_Unicode(t *testing.T) {
	input := "Привет мир"
	want := "Привет_мир"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_ExceedingMaxLength(t *testing.T) {
	input := ""
	for range 150 {
		input += "a"
	}
	want := ""
	for range 100 {
		want += "a"
	}
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(long) = %q; want %q", got, want)
	}
}

func TestSanitizeTitle_MixedValidInvalidUnicode(t *testing.T) {
	input := "  Résumé: John/Smith?* "
	want := "Résumé_John_Smith"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_OnlyUnderscores(t *testing.T) {
	input := "____"
	want := "untitled"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_Emoji(t *testing.T) {
	input := "Project 🚀 Launch"
	want := "Project_🚀_Launch"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}

func TestSanitizeTitle_NewlineTab(t *testing.T) {
	input := "foo\nbar\tbaz"
	want := "foo_bar_baz"
	got := sanitizeTitle(input)
	if got != want {
		t.Errorf("sanitizeTitle(%q) = %q; want %q", input, got, want)
	}
}
