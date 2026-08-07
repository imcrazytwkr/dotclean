package merge

import "testing"

func TestSkipAttr(t *testing.T) {
	if !skipAttr(xattrMacl, true) || !skipAttr(xattrProvenance, false) {
		t.Fatal("macl/provenance must always skip")
	}
	if !skipAttr(xattrQuarantine, false) {
		t.Fatal("quarantine skipped without SetQuarantine")
	}
	if skipAttr(xattrQuarantine, true) {
		t.Fatal("quarantine applied with SetQuarantine")
	}
	if skipAttr("com.apple.lastuseddate#PS", false) {
		t.Fatal("lastuseddate must not skip")
	}
}
