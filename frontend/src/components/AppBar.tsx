import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { logoutUser } from "../features/authSlice";
import { AppBar, Box, Button, IconButton, Toolbar, Typography } from "@mui/material";
import { AccountCircle, Brightness4, Brightness7, Forum, Logout, Person } from "@mui/icons-material";
import { useTheme } from "../contexts/ThemeContext";

export default function ForumAppBar() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { token, userID } = useAppSelector(state => state.auth);
    const { darkMode, toggleDarkMode } = useTheme();

    const handleLogout = () => {
        dispatch(logoutUser());
        navigate('/login');
    }

    return (
        <AppBar position="sticky">
            <Toolbar>
                <Forum sx={{ mr: 2 }} />
                <Typography
                    variant="h6"
                    component="div"
                    sx={{
                        flexGrow: 1,
                        cursor: "pointer"
                    }}
                    onClick={() => navigate(token ? '/topics' : '/login')}
                >
                    Gossip with Go
                </Typography>

                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1,
                    }}
                >
                    {/* Theme Toggle */}
                    <IconButton 
                        color="inherit"
                        onClick={toggleDarkMode}
                        title={`Switch to ${darkMode ? 'Light' : 'Dark'} Mode`}
                    >
                        {darkMode ? <Brightness7 /> : <Brightness4 />}
                    </IconButton>
                </Box>

                {token 
                    ? (
                        <>
                            <IconButton
                                color="inherit"
                                onClick={() => navigate(`/users/${userID}`)}
                                title="My Profile"
                            >
                                <AccountCircle />
                            </IconButton>

                            <Button
                                color="inherit"
                                startIcon={<Logout />}
                                onClick={handleLogout}
                            >
                                Logout
                            </Button>
                        </>
                    )
                    : (
                        <>
                            <Button
                                color="inherit"
                                startIcon={<Person />}
                                onClick={() => navigate('/login')}
                            >
                                Login
                            </Button>

                            <Button
                                color="inherit"
                                onClick={() => navigate('/register')}
                            >
                                Register
                            </Button>
                        </>
                    )
                }
            </Toolbar>
        </AppBar>
    );
}