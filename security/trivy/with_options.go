package main

// WithVersion définit la version de l'image Trivy (pattern immuable).
func (m *Trivy) WithVersion(version string) *Trivy {
	return &Trivy{
		Version:          version,
		Severity:         m.Severity,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}

// WithSeverity définit les sévérités à signaler, ex: "CRITICAL" ou "HIGH,CRITICAL".
func (m *Trivy) WithSeverity(severity string) *Trivy {
	return &Trivy{
		Version:          m.Version,
		Severity:         severity,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}
