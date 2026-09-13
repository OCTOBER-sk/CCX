package migration

// Config migration fixtures per spec 41 (lines 1410-1427) and spec 150 (lines 4748-4771):
// Every migration: preserve old config, validate new schema, create backup, report changes, reversible.
// Migration fixtures must be kept permanently.

type Fixture struct {
	SchemaVersion string
	BackupCreated bool
	Reversible bool
}

func LoadFixture(version string) Fixture {
	return Fixture{
		SchemaVersion: version,
		BackupCreated: true,
		Reversible: true,
	}
}
