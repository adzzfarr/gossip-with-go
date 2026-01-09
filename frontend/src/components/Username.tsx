import { Typography, type TypographyProps } from "@mui/material";
import type { MouseEvent } from "react";
import { useNavigate } from "react-router-dom";

interface UsernameProps extends Omit<TypographyProps, 'onClick' | 'prefix'> {
    username: string;
    userID: number;
    prefix?: boolean; // Show '/u' before username
}

export default function Username({ username, userID, prefix = true, ...props }: UsernameProps) {
    const navigate = useNavigate();

    const handleClick = (e: MouseEvent) => {
        e.stopPropagation();
        navigate(`/users/${userID}`);
    }

    const displayName = prefix ? `u/${username}` : username;

    return (
        <Typography
            {...props}
            onClick={handleClick}
            sx={{
                cursor: 'pointer',
                '&:hover': {
                    textDecoration: 'underline',
                    color: 'primary.main',
                },
                transition: 'color 0.2s',
                ...props.sx,
            }}
        >
            {displayName}
        </Typography>
    );
}