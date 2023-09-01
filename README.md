# ObservInstaller installer

Observinstaller is a tool which can be helpful in setting up the native environment (Bare-metal, VM) with observability. It comes with various options which can help ease the job of the user to manage the observability appliction at some level

---
### If you run it with --help you're presented with this:
    USAGE: observinstaller <command> <options>
    Supported commands:
      download
      run
      kill
      remove
      otel

The application works based on a .config.yml file that should be present alongside the binary which is read by the application to do the work. The .config.yml has the description for the packages, installation, download directory set-up and base default config layout for building a OpenTelemetry collector config file which can be modified as per the user's needs.

---
### Specs of the .config.yml:
* *downloadDirectory*: Defines where the downloaded things should go.
* *installationDirectory*: Defines where should the installed items go, it is generally preferred to keep this item inside the user's directory for easy cleanup and so as to not mess with system entities in any way
* *baseOtelConfig*: Describes what default config should generated otel config file must have, it's identical to Otel config's service section
* *packages*: It's a list of the items which contains the configurations for installation, running and otel config for overriding the default one and installationMode support which is either full or minimal, these two options generally governs in what type of installation this package can be used. This will be more understandable if we understand the tool a bit more

---
### ObservInstaller options:
* *download*: Download helps the user to download and configure the applications based on what mode is selected. observeinstaller has 2 modes -minimal and -full. So if a package is marked with one of these and you select that option to install the package will be selected for installation. 
	* For eg: If in .config.yml let's say Grafana is marked with minimal (in installModeSupport) and you run the download option with -minimal then Grafana will be downloaded and installed if you've marked Grafana for both minimal and full installation, then it'll be installed in both selections
* *run*: Run option is used to run the installed application based on the selected modes -minimal or -full. Application is started in the background and the user is notified of the result
* *kill*: Does the opposite of run, it kills the running application. It also has an option to restart the running application if in case, it is required. This option generally works based on the data saved by the run command
* *remove*: This option can be used to remove the either the installation directory or the downlaods directory or both of them. This option can be considered as an uninstallation step
* *otel*: This option works based on the baseOtelConfig field present in .config.yml it takes those fields and generate a new yaml `otel-config.yaml` for the otel collector. The fields are default, so people can change according to their needs. These configs can be overwritten using the field (pkgOtelConfig) in the package

---
