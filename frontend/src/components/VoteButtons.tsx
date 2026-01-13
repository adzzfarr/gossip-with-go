import { useNavigate } from "react-router-dom";
import { useAppSelector } from "../hooks/redux";
import { useState } from "react";
import { Box, CircularProgress, IconButton, Tooltip, Typography } from "@mui/material";
import { ArrowDownward, ArrowUpward } from "@mui/icons-material";
import { apiClient } from "../api/client";

interface VoteButtonsProps {
    postID?: number;
    commentID?: number;
    initialVoteCount?: number; 
    initialUserVote?: 1 | -1 | null;
    size?: 'small' | 'medium' | 'large';
    orientation?: 'vertical' | 'horizontal';
    onVoteChange?: (data: {
        voteCount: number;
        userVote: 1 | -1 | null;
    }) => void;
}

export default function VoteButtons({
    postID,
    commentID,
    initialVoteCount = 0,
    initialUserVote = null,
    size = 'medium', // 'small' | 'medium' | 'large'
    orientation = 'vertical', // 'vertical' | 'horizontal'
    onVoteChange,
}: VoteButtonsProps) {
    const navigate = useNavigate();

    const { isAuthenticated } = useAppSelector(state => state.auth);

    const [voteCount, setVoteCount] = useState(initialVoteCount);
    const [userVote, setUserVote] = useState<1 | -1 | null>(initialUserVote);
    const [isLoading, setIsLoading] = useState(false);

    const handleVote = async (voteType: 1 | -1) => {
        // Check authentication
        if (!isAuthenticated) {
            navigate('/login');
            return;
        }

        // Prevent multiple votes at the same time
        if (isLoading) return;

        const previousVoteCount = voteCount;
        const previousUserVote = userVote;

        // Determine new vote state
        let newVoteCount = voteCount;
        let newUserVote: 1 | -1 | null = userVote;

        if (userVote === voteType) {
            // Toggle off (remove vote)
            newVoteCount = voteCount - voteType;
            newUserVote = null;
        } else if (userVote !== null) {
            // Switch vote
            newVoteCount = voteCount - userVote + voteType;
            newUserVote = voteType;
        } else {
            // New vote
            newVoteCount = voteCount + voteType;
            newUserVote = voteType;
        }

        setVoteCount(newVoteCount);
        setUserVote(newUserVote);
        setIsLoading(true);

        try {
            const endpoint = postID
                ? `posts/${postID}/vote`
                : `comments/${commentID}/vote`

            const response = await apiClient.post(
                endpoint,
                { voteType }
            );
            const data = response.data;

            // Update state based on server response
            setVoteCount(data.voteCount);
            setUserVote(data.userVote);

            // Notify parent component of vote change
            if (onVoteChange) {
                onVoteChange({
                    voteCount: data.voteCount,
                    userVote: data.userVote,
                });
            }
        } catch (error) {
            // Revert to previous state on error
            setVoteCount(previousVoteCount);
            setUserVote(previousUserVote);
            console.error('Error submitting vote:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const getIconSize = () => {
        switch (size) {
            case 'small':
                return 'small';
            case 'large':
                return 'large';
            default:
                return 'medium';
        }
    };

    const getVoteColor = (voteType: 1 | -1) => {
        if (userVote === voteType) {
            return voteType === 1 ? 'warning' : 'info';
        }
        return 'default';
    }
 
    const isVertical = orientation === 'vertical';

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: isVertical ? 'column' : 'row',
                alignItems: 'center',
                gap: isVertical ? 0 : 1,
            }}
        >
            {/* Upvote */}
            <Tooltip
                title={isAuthenticated ? 'Upvote' : 'Login to vote'}
                arrow 
            > 
                <span>
                    <IconButton
                        size={getIconSize()}
                        color={getVoteColor(1)}
                        onClick={() => handleVote(1)}
                        disabled={isLoading || !isAuthenticated}
                        sx={{
                            '&:hover': {
                                backgroundColor: 'rgba(255, 165, 0, 0.1)'
                            },
                            transition: 'all 0.2s',
                        }}
                    >
                        {
                            isLoading && userVote === 1
                                ? (<CircularProgress size={24} />)
                                : (<ArrowUpward />)
                        }
                    </IconButton>
                </span>
            </Tooltip>

            {/* Vote Count */}
            <Typography
                variant={size === 'small' ? 'body2' : 'body1'}
                sx={{
                    fontWeight: userVote ? 'bold' : 'normal',
                    color: userVote === 1
                        ? 'warning.main'
                        : userVote === -1
                            ? 'info.main'
                            : 'text.primary',
                    minWidth: size === 'small' ? 24 : 32,
                    textAlign: 'center',
                    transition: 'all 0.2s',
                }}
            >
                {voteCount}
            </Typography>

            {/* Downvote */}
            <Tooltip
                title={isAuthenticated ? 'Downvote' : 'Login to vote'}
                arrow 
            >
                <span>
                    <IconButton
                        size={getIconSize()}
                        color={getVoteColor(-1)}
                        onClick={() => handleVote(-1)}
                        disabled={isLoading || !isAuthenticated}
                        sx={{
                            '&:hover': {
                                backgroundColor: 'rgba(0, 191, 255, 0.1)'
                            },
                            transition: 'all 0.2s',
                        }}
                    >
                        {
                            isLoading && userVote === -1
                                ? (<CircularProgress size={24} />)
                                : (<ArrowDownward />)
                        }
                    </IconButton>
                </span>
            </Tooltip>
        </Box>
    )
}