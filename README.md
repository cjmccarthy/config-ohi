New Relic Infrastructure - Application Config Ingestor
======================================================

This is an on-host integration for New Relic Infrastructure. Its primary purpose is to ingest application configurations into New Relic as inventory data, thought it can be used to ingest any file.

## Getting Started

### Compilation
To begin, read the official New Relic documentation on what an infrastructure integration is, and how to get started with the SDK

[Integration Description](https://docs.newrelic.com/docs/integrations/integrations-sdk/getting-started/intro-infrastructure-integrations-sdk)

[SDK Setup](https://github.com/newrelic/infra-integrations-sdk)

To download the source to your GOPATH:

`$ go get github.com/cjmccarthy/config-ohi`

Then, in that directory, build the executable. If you are cross-compiling, you will need to add your target OS and architecture at this stage. For instance, if you are cross compiling from OSX to Linux running on x86-64:

`env GOOS=linux GOARCH=amd64 go build`

### Config and Definition Files

In the source, you will find two included files:
```
config-ohi-definition.yml
config-ohi-config.yml
```

These files are standard required configurations for the SDK. Their deployment will be covered in the next section.

You will see that the config-ohi-config.yml file has, on line 6:

`nr_ingest_conf: ./test.yml`

The test.yml file included in this repository contains examples of how to define which files to send to New Relic. Using that as a guide, you should write your own file. Then, replace that line in config-ohi-config.yml with your own filename and path if you so choose. Each item in the file should contain:

inventory: The path where this item will be located in New Relic.

path: The actual path on the system where the file is located.

type: The file type. Currently there are four options described below.

| type | Behavior |
--- | --- |
text | Imported as-is|
yaml | Top level key-value pairs will be unwrapped. Any nested data will appear as json |
json |  Top level key-value pairs will be unwrapped. Any nested data will appear as json |
xml | Top level key-value pairs will be unwrapped. Any nested data will appear as json. Since xml is order-dependent, sequence numbers will be added alongside the original data |
~~properties~~ | coming soon! |

If your configurations are highly nested, it is recommended to use the "text" type for readability. However, if they are fairly flat, picking the correct type will make the inventory data more rich by breaking down changes to the key-value level.

The final step is to define the prefix under which you would like these inventory items to appear inside of the Infrastructure UI. To do this, modify line 10 of config-ohi-definition.yml

`prefix: my/prefix`

Whatever prefix is placed here will be prepended to the "inventory" paths given in the file containing information about where your files are stored. For instance, test.yml's "inventory" paths are in the form `my/path/one`. So, if left as is, the paths would appear in inventory as:

```
my/prefix/my/path/one
my/prefix/my/path/two
...
```

### Deployment

Once the files are properly configured, the standard deployment instructions for the executable and the two files described above can be found [here](https://docs.newrelic.com/docs/integrations/integrations-sdk/getting-started/integration-file-structure-activation)

The file that you have created in the last section to replace test.yml should be deployed at the path you have specified. If you have not specified a new path, it should be placed alongside the executable.

### Data

Once the integration is running, you should see your config files added in the Inventory tab in infrastructure under the prefix and path you have specified. Changes to these files will now show up in the event stream with a before-after diff.


### Debugging 

https://docs.newrelic.com/docs/integrations/integrations-sdk/troubleshooting/not-seeing-infrastructure-integration-data

