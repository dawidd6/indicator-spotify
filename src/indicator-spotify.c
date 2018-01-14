#include <libappindicator/app-indicator.h>

#define INTERVAL 3
#define MAX_WIDTH 50

GtkWidget *menu;
GtkWidget *item_quit;
AppIndicator *indicator;
GDBusConnection *bus;
gchar bundle[MAX_WIDTH];
GVariant *result, *props;
gchar **artists = NULL,
		*artist = NULL,
		*title = NULL;

gboolean update()
{
	result = g_dbus_connection_call_sync(bus,
			"org.mpris.MediaPlayer2.spotify",
			"/org/mpris/MediaPlayer2",
			"org.freedesktop.DBus.Properties",
			"Get",
			g_variant_new("(ss)",
			"org.mpris.MediaPlayer2.Player",
			"Metadata"),
			G_VARIANT_TYPE("(v)"),
			G_DBUS_CALL_FLAGS_NONE, -1, NULL, NULL);

	if(!result)
	{
		if(app_indicator_get_status(indicator) != APP_INDICATOR_STATUS_PASSIVE)
			app_indicator_set_status(indicator, APP_INDICATOR_STATUS_PASSIVE);
		return 1;
	}

	g_variant_get(result, "(v)", &props);
	g_variant_lookup(props, "xesam:artist", "^as", &artists);
	g_variant_lookup(props, "xesam:title", "s", &title);

	if(artists)
		artist = g_strjoinv(", ", artists);
	else
		artist = "(Unknown Artist)";

	if(!title)
		title = "(Unknown Song)";

	snprintf(bundle, MAX_WIDTH, "%s – %s", artist, title);

	if(app_indicator_get_status(indicator) != APP_INDICATOR_STATUS_ACTIVE)
		app_indicator_set_status(indicator, APP_INDICATOR_STATUS_ACTIVE);

	app_indicator_set_label(indicator, bundle, NULL);

	return 1;
}

int main(int argc, char *argv[])
{
	gtk_init(&argc, &argv);

	menu = gtk_menu_new();
	indicator = app_indicator_new("indicator-spotify", "renamed-spotify-client", APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
	app_indicator_set_menu(indicator, GTK_MENU(menu));
	app_indicator_set_title(indicator, "Spotify Indicator");

	item_quit = gtk_menu_item_new_with_label("Quit");
	gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_quit);
	g_signal_connect(item_quit, "activate", gtk_main_quit, NULL);

	bus = g_bus_get_sync(G_BUS_TYPE_SESSION, NULL, NULL);
	if(!bus)
		return 1;

	g_timeout_add_seconds(INTERVAL, update, NULL);

	gtk_widget_show_all(menu);
	gtk_main();
	g_object_unref(bus);
	return 0;
}
