#Assembled by dawidd6
PROGRAM=indicator-spotify
START_COLOR=\033[0;33m
CLOSE_COLOR=\033[m
DESTDIR=

build:
	@go build -tags gtk_3_18

install:
	@install -d $(DESTDIR)/usr/bin
	@install -d $(DESTDIR)/usr/share/applications
	@install -d $(DESTDIR)/etc/xdg/autostart
	@install -d $(DESTDIR)/usr/share/icons/hicolor/scalable/apps
	@install $(PROGRAM) $(DESTDIR)/usr/bin
	@install -m 644 $(PROGRAM).desktop $(DESTDIR)/usr/share/applications
	@install -m 644 $(PROGRAM).desktop $(DESTDIR)/etc/xdg/autostart
	@install -m 644 $(PROGRAM).svg $(DESTDIR)/usr/share/icons/hicolor/scalable/apps

uninstall:
	@rm -rf /usr/bin/$(PROGRAM)
	@rm -rf /usr/share/applications/$(PROGRAM).desktop
	@rm -rf /etc/xdg/autostart/$(PROGRAM).desktop
	@rm -rf /usr/share/icons/hicolor/scalable/apps/$(PROGRAM)

.PHONY: install uninstall
