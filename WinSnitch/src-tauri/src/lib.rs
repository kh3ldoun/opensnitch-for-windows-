use tauri::{menu::{MenuBuilder, MenuItemBuilder}, tray::TrayIconBuilder, Manager};

#[tauri::command]
fn service_status() -> &'static str {
    "running"
}

pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_autostart::init(
            tauri_plugin_autostart::MacosLauncher::LaunchAgent,
            None,
        ))
        .invoke_handler(tauri::generate_handler![service_status])
        .setup(|app| {
            let show = MenuItemBuilder::new("Open WinSnitch").id("open").build(app)?;
            let quit = MenuItemBuilder::new("Quit").id("quit").build(app)?;
            let menu = MenuBuilder::new(app).items(&[&show, &quit]).build()?;
            TrayIconBuilder::new().menu(&menu).on_menu_event(|app, event| match event.id().as_ref() {
                "open" => {
                    if let Some(win) = app.get_webview_window("main") { let _ = win.show(); let _ = win.set_focus(); }
                }
                "quit" => app.exit(0),
                _ => {}
            }).build(app)?;
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
