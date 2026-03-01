# OpenSnitch-Windows KMDF Callout Driver

This directory contains the kernel-mode WFP callout driver skeleton for pended connection decisions.

Planned components:
- KMDF driver project (`OpenSnitchCallout.sys`)
- ALE connect/recv accept callout registrations
- IOCTL interface for decision handoff with daemon
- Optional test-signing support for developer mode
