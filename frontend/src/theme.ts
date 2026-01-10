import { createTheme, type Theme } from '@mui/material/styles';

// Shared Component Styling
const sharedComponents = {
    MuiButton: {
        styleOverrides: {
            root: {
                textTransform: 'none',
                borderRadius: '8px',
            }
        }
    },
    MuiCard: {
        styleOverrides: {
            root: {
                borderRadius: '12px',
            }
        }
    }
}

// Shared Typography 
const sharedTypography = {
    fontFamily: `"Roboto", "Helvetica", "Arial", sans-serif`,
    h1: {
        fontSize: '2.5rem',
        fontWeight: 600,
    },
    h2: {
        fontSize: '2rem',
        fontWeight: 600,
    },
    h3: {
        fontSize: '1.75rem',
        fontWeight: 500,
    },
    h4: {
        fontSize: '1.5rem',
        fontWeight: 500,
    },
    h5: {
        fontSize: '1.25rem',
        fontWeight: 500,
    },
    h6: {
        fontSize: '1rem',
        fontWeight: 500,
    },
};

export const lightTheme: Theme = createTheme({
    palette: {
        mode: 'light',
        primary: {
            main: '#1976d2',
            light: '#63a4ff',
            dark: '#004ba0',
        },
        secondary: {
            main: '#9c27b0',
            light: '#d05ce3',
            dark: '#6a0080',
        },
        background: {
            default: '#f5f5f5',
            paper: '#ffffff',
        },
    },
    typography: sharedTypography,
    components: {
        ...sharedComponents,
        MuiCard: {
            styleOverrides: {
                root: {
                    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
                    borderRadius: '12px',
                }
            }
        }
    }
});

export const darkTheme: Theme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#90caf9',
            light: '#e3f2fd',
            dark: '#42a5f5',
        },
        secondary: {
            main: '#ce93d8',
            light: '#f3e5f5',
            dark: '#ab47bc',
        },
        background: {
            default: '#121212',
            paper: '#1e1e1e',
        },
    },
    typography: sharedTypography,
    components: {
        ...sharedComponents,
        MuiCard: {
            styleOverrides: {
                root: {
                    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.5)', // Darker shadow
                    borderRadius: '12px',
                }
            }
        },
    }
});