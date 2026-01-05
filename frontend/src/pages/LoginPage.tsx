import { useEffect, useState, type ChangeEvent, type FormEvent } from "react";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useNavigate, Link as RouterLink } from "react-router-dom";
import { clearError, loginUser } from "../features/authSlice";
import { Alert, Box, Button, CircularProgress, Container, Link, Paper, TextField, Typography } from "@mui/material";

export default function LoginPage() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { loading, error, isAuthenticated } = useAppSelector(state => state.auth); 

    // Clear error when component mounts
    useEffect(() => {
        return () => {
            dispatch(clearError());
        };
    }, []);
    
    // Redirect if authenticated
    useEffect(() => {
        if (isAuthenticated) {
            navigate('/');
        }
    }, [isAuthenticated, navigate]);

    // Clear errors when user starts typing
    function handleUsernameChange(e: ChangeEvent<HTMLInputElement>) {
        setUsername(e.target.value);

        if (error) {
            dispatch(clearError());
        }
    }

    function handlePasswordChange(e: ChangeEvent<HTMLInputElement>) {
        setPassword(e.target.value);

        if (error) {
            dispatch(clearError());
        }
    }

    // Handle form submission
    function handleSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

        // Send login request to backend via Redux thunk
        dispatch(loginUser({ username, password }));
    }

    return (
        <Container component='main' maxWidth='xs'>
            <Box
                sx={{
                    mt: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Paper
                    elevation={3}
                    sx={{
                        p: 4,
                        width: '100%',
                    }}
                >
                    <Typography
                        component='h1'
                        variant='h4'
                        align='center'
                        gutterBottom
                    >
                        Login
                    </Typography>

                    {error && (
                        <Alert severity='error' sx={{ mb: 2 }}>
                            {error}
                        </Alert>
                    )}

                    <Box
                        component='form'
                        onSubmit={handleSubmit}
                        sx={{ mt: 1}}
                    >
                        <TextField 
                            required
                            id='username'
                            label='Username'
                            name='username'
                            value={username}
                            onChange={handleUsernameChange}
                            disabled={loading}
                            autoComplete='username'
                            margin='normal'
                            fullWidth
                            autoFocus
                        />

                        <TextField 
                            required
                            id='password'
                            label='Password'
                            name='password'
                            type='password'
                            value={password}
                            onChange={handlePasswordChange}
                            disabled={loading}
                            autoComplete='current-password'
                            margin='normal'
                            fullWidth
                        />

                        <Button
                            type='submit'
                            variant='contained'
                            disabled={loading || !username || !password}
                            fullWidth
                            sx={{ mt: 3, mb: 2 }}
                        >
                            {loading ? <CircularProgress size={24} /> : 'Login'}
                        </Button>

                        <Box sx={{ textAlign: 'center' }}>
                            <Typography variant='body2'>
                                Don't have an account?{' '}
                                <Link
                                    component={RouterLink}
                                    to='/register'
                                    onClick={() => { // Clear error on navigation
                                        dispatch(clearError());
                                    }}
                                >
                                    Register
                                </Link>
                            </Typography>
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
}