package main

import (
	"fmt"


	"golang.org/x/sys/windows"
)

// getProcessPath returns the full path of the executable for a given PID.
func getProcessPath(pid uint32) (string, error) {
	// Require PROCESS_QUERY_LIMITED_INFORMATION or PROCESS_QUERY_INFORMATION
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("failed to open process %d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	var path [windows.MAX_PATH]uint16
	size := uint32(windows.MAX_PATH)

	// QueryFullProcessImageNameW
	err = windows.QueryFullProcessImageName(handle, 0, &path[0], &size)
	if err != nil {
		return "", fmt.Errorf("QueryFullProcessImageNameW failed: %w", err)
	}

	return windows.UTF16ToString(path[:size]), nil
}

// getProcessUser returns the SID string of the user running the process.
func getProcessUser(pid uint32) (string, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("failed to open process %d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	var token windows.Token
	err = windows.OpenProcessToken(handle, windows.TOKEN_QUERY, &token)
	if err != nil {
		return "", fmt.Errorf("failed to open process token: %w", err)
	}
	defer token.Close()

	user, err := token.GetTokenUser()
	if err != nil {
		return "", fmt.Errorf("failed to get token user: %w", err)
	}

	return user.User.Sid.String(), nil
}

// getCommandLine returns the command line string for a given PID.
// Extracting command line from another process requires reading its PEB using NtQueryInformationProcess.
func getCommandLine(pid uint32) (string, error) {
	// Requires PROCESS_VM_READ and PROCESS_QUERY_INFORMATION
	handle, err := windows.OpenProcess(windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("failed to open process %d for VM_READ: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	// In a real production system, extracting command lines across x86/x64 requires
	// dealing with Wow64 and reading PEB correctly.
	// For MVP of OpenSnitch-Windows, we map this out to get basic path + user first.

	return "<Command Line Not Available>", nil
}
