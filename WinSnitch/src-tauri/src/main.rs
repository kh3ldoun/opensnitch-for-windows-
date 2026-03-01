#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use tauri::{AppHandle, Manager};

// A simple command to ping the Go daemon (placeholder)
#[tauri::command]
async fn ping_backend(app: AppHandle) -> Result<String, String> {
    // In a real implementation, this would communicate with the Go backend via RPC/Sockets/NamedPipe
    Ok("Pong from WinSnitch Backend".to_string())
}

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![ping_backend])
        .setup(|app| {
            // Setup Tray Icon, Notification Windows, etc.
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
