import { useEffect, useState, type ChangeEvent, type FormEvent } from "react";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { clearError, registerUser } from "../features/auth/authSlice";
import { Alert, Box, Button, CircularProgress, Container, Link, List, ListItem, ListItemIcon, ListItemText, Paper, TextField, Typography } from "@mui/material";
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';

export default function RegisterPage() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [validationError, setValidationError] = useState('');

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

    // Password Validation Helpers (match backend constraint)
    const hasMinLength = password.length >= 8;
    const hasLowercase = /[a-z]/.test(password);
    const hasUppercase = /[A-Z]/.test(password);
    const hasDigit = /[0-9]/.test(password);
    const isValidPassword = hasMinLength && hasLowercase && hasUppercase && hasDigit;
    const passwordsMatch = password === confirmPassword && confirmPassword !== '';

    // Clear errors when user starts typing
    function handleUsernameChange(e: ChangeEvent<HTMLInputElement>) {
        setUsername(e.target.value);
        setValidationError('');

        if (error) {
            dispatch(clearError());
        }
    }

    function handlePasswordChange(e: ChangeEvent<HTMLInputElement>) {
        setPassword(e.target.value);
        setValidationError('');

        if (error) {
            dispatch(clearError());
        }
    }

    function handleConfirmPasswordChange(e: ChangeEvent<HTMLInputElement>) {
        setConfirmPassword(e.target.value);
        setValidationError('');

        // No need to dispatch clearError (no backend error, only local validation)
    }

    // Handle form submission
    function handleSubmit(e: FormEvent<HTMLFormElement>) { 
        e.preventDefault();

        // Client-side Validation
        // Passwords match
        if (password !== confirmPassword) {
            setValidationError('Passwords do not match.');
            return;
        }

        // Password passes backend constraints
        if (!hasMinLength) {
            setValidationError('Password must be at least 8 characters long.');
            return;
        }

        if (!hasLowercase) {
            setValidationError('Password must contain at least one lowercase letter.');
            return;
        }

        if (!hasUppercase) {
            setValidationError('Password must contain at least one uppercase letter.');
            return;
        }

        if (!hasDigit) {
            setValidationError('Password must contain at least one digit.');
            return;
        }

        // Send register request to backend via Redux thunk if validation passes
        dispatch(registerUser({ username, password }) );
    };

    const displayError = error || validationError;

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
                        Register
                    </Typography>

                    {displayError && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {displayError}
                        </Alert>
                    )}

                    <Box 
                        component="form" 
                        onSubmit={handleSubmit} 
                        sx={{ mt: 1 }}
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
                            autoComplete='new-password'
                            margin='normal'
                            fullWidth
                        />

                        {password && (
                            <Box sx={{ mt: 1, mb: 2 }}>
                                <Typography
                                    variant='caption'
                                    color='text.secondary'
                                    gutterBottom
                                >
                                    Password Requirements:
                                </Typography>
                                <List dense>
                                    <ListItem disableGutters>
                                        <ListItemIcon sx={{ minWidth: 32 }}>
                                            {
                                                hasMinLength
                                                    ? (<CheckCircleIcon color="success" fontSize="small" />)
                                                    : (<CancelIcon color="error" fontSize="small" />) 
                                            }
                                        </ListItemIcon>
                                        <ListItemText 
                                            primary="At least 8 characters long"
                                            slotProps={{ primary: { variant: 'body2' } }}
                                        />
                                    </ListItem>
                                    <ListItem disableGutters>
                                        <ListItemIcon sx={{ minWidth: 32 }}>
                                            {
                                                hasLowercase
                                                    ? (<CheckCircleIcon color="success" fontSize="small" />)
                                                    : (<CancelIcon color="error" fontSize="small" />) 
                                            }
                                        </ListItemIcon>
                                        <ListItemText 
                                            primary="At least one lowercase letter"
                                            slotProps={{ primary: { variant: 'body2' } }}
                                        />
                                    </ListItem>
                                    <ListItem disableGutters>
                                        <ListItemIcon sx={{ minWidth: 32 }}>
                                            {
                                                hasUppercase
                                                    ? (<CheckCircleIcon color="success" fontSize="small" />)
                                                    : (<CancelIcon color="error" fontSize="small" />) 
                                            }
                                        </ListItemIcon>
                                        <ListItemText 
                                            primary="At least one uppercase letter"
                                            slotProps={{ primary: { variant: 'body2' } }}
                                        />
                                    </ListItem>
                                    <ListItem disableGutters>
                                        <ListItemIcon sx={{ minWidth: 32 }}>
                                            {
                                                hasDigit
                                                    ? (<CheckCircleIcon color="success" fontSize="small" />)
                                                    : (<CancelIcon color="error" fontSize="small" />) 
                                            }
                                        </ListItemIcon>
                                        <ListItemText 
                                            primary="At least one digit"
                                            slotProps={{ primary: { variant: 'body2' } }}
                                        />
                                    </ListItem>
                                </List>
                            </Box>
                        )}

                        <TextField 
                            required
                            id='confirmPassword'
                            label='Confirm Password'
                            name='confirmPassword'
                            type='password'
                            value={confirmPassword}
                            onChange={handleConfirmPasswordChange}
                            disabled={loading}
                            autoComplete='new-password'
                            margin='normal'
                            fullWidth
                            error={confirmPassword !== '' && !passwordsMatch}
                            helperText={
                                confirmPassword !== '' && !passwordsMatch
                                    ? 'Passwords do not match.'
                                    : ''
                            }
                        />

                        <Button
                            type='submit'
                            variant='contained'
                            disabled={loading || !username || !password || !confirmPassword || !isValidPassword || !passwordsMatch}
                            fullWidth
                            sx={{
                                mt: 3,
                                mb: 2,
                            }}
                        >
                            {loading ? <CircularProgress size={24} /> : 'Register'}
                        </Button>

                        <Box sx={{ textAlign: 'center' }}>
                            <Typography variant='body2'>
                                Already have an account?{' '}
                                <Link 
                                    component={RouterLink} 
                                    to="/login"
                                    onClick={() => { // Clear error on navigation
                                        dispatch(clearError());
                                    }}
                                >
                                    Login
                                </Link>
                            </Typography>
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    )
}