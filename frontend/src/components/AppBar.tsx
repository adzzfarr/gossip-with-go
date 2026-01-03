import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { logoutUser } from "../features/auth/authSlice";
import { AppBar, Box, Button, Toolbar, Typography } from "@mui/material";
import { Forum, Logout } from "@mui/icons-material";

export default function ForumAppBar() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { username } = useAppSelector(state => state.auth);

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
                        flexGrow: 0,
                        cursor: "pointer"
                    }}
                    onClick={() => navigate('/topics')}
                >
                    Gossip with Go
                </Typography>

                <Box sx={{ flexGrow: 1 }} />

                <Typography variant="body1" sx={{ mr: 2 }}>
                    {username}
                </Typography>

                <Button 
                    onClick={handleLogout}
                    color="inherit"
                    startIcon={<Logout />}
                >
                    Logout
                </Button>
            </Toolbar>
        </AppBar>
    );
}