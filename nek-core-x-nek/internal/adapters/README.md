# Transport adapter plan

A production adapter should validate protocol config, render engine config without logging secrets, validate before reload, apply changes atomically, reconcile desired users, and report health/usage. Suggested order: Xray-core for VLESS + REALITY, then other supported Xray protocols; separate Hysteria 2 and TUIC engines; WireGuard via tightly controlled OS tooling. The current protocol catalog is informational only.
