# Calculator

Calculator that can print a file to the Default Network Printer containing the last calculation made.

**Currently only Windows is supported.**

## How to install the application

To install the application, you can simply go to the [Releases](https://github.com/Feggah/calculator/releases) section of this repository and download the desired version or [click here](https://github.com/Feggah/calculator/releases/download/v0.1.0/calculator.exe) to download the latest version.

## Demo

The application home screen looks like a standard calculator:

![image](/img/home-screen.png)

Below is a calculation example. Please note that the application will only save the last calculation, so if you do a single inline equation, it would be  saved entirely.

If you do multiple calculations to reach the result, only the last one would be saved. In the example below, it would save `90 + 50` instead of the whole equation.

![image](/img/calculation-example.png)

![image](/img/calculation-made-example.png)

Now that the equation is saved (you clicked `=`), you can print the equation by going to the menu `Arquivo` and then `Imprimir`, or simply use the shortcut `CTRL + P`.

![image](/img/print-title-example.png)

In the printscreen above you should type the title of the calculation and press `Confirmar` to print it or `Cancelar` to cancel.

After pressing `Confirmar`, the content that would be printed is going to be the same as the printscreen below:

![image](/img/printed-content-example.png)

## License

Calculator is under the Apache 2.0 license.

## Contributing

Feel free to open issues or pull requests to this repository. Contributions are more than welcome.
