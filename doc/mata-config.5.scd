mata-config(5)

# Name

mata-config - configuration file formats for *mata*(1)

# Configuration

There is only one configuration file for mata: *config.json*. The program looks 
after this file in your XDG config home plus "mata", which defaults to 
~/.config/mata.

This file uses the _json_ format.

# CONFIG.JSON

This file is used to configure the behavior of mata.

## OPTIONS

*endpoint*
	Sets the endpoint that will be used to send resquests.

	Default: https://mataroa.blog/api

*key*
	The API key provided by your mataroa dashboard.

# SEE ALSO

*mata*(1)

# AUTHORS

Created by Victor Freire <victor@freire.dev.br>. For more information about mata
development, see https://sr.ht/~glorifiedgluer/mata.