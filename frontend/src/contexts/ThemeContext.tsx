import { createContext, useContext, useState, type ReactNode } from "react";
import { CssBaseline, ThemeProvider as MuiThemeProvider } from "@mui/material";
import { darkTheme, lightTheme } from "../theme";

interface ThemeContextType {
    darkMode: boolean;
    toggleDarkMode: () => void;
}

const ThemeContext = createContext<ThemeContextType>({
    darkMode: true,
    toggleDarkMode: () => {},
});

export const useTheme = () =>  useContext(ThemeContext);

interface ThemeProviderProps {
    children: ReactNode
}

export function ThemeProvider({ children }: ThemeProviderProps) {
    const [darkMode, setDarkMode] = useState<boolean>(
        () => {
            const storedPref = localStorage.getItem('darkMode');
            return storedPref === 'true';
        }
    );

    const toggleDarkMode = () => {
        setDarkMode((prevMode) => {
            const newMode = !prevMode;
            localStorage.setItem('darkMode', newMode.toString());
            return newMode;
        });
    };

    return (
        <ThemeContext.Provider value={{ darkMode, toggleDarkMode }}>
            <MuiThemeProvider theme={darkMode ? darkTheme : lightTheme}>
                <CssBaseline />
                {children}
            </MuiThemeProvider>
        </ThemeContext.Provider>
    );
}