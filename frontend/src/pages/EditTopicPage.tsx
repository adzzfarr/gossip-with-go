import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState, type ChangeEvent, type FormEvent } from "react";
import { clearError, updateTopic } from "../features/topicsSlice";
import { Alert, Box, Button, CircularProgress, Container, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";

export default function EditTopicPage() {
    const { topicID} = useParams<{ topicID: string }>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { topics, loading, submitting, submitError } = useAppSelector(state => state.topics);
    const { userID } = useAppSelector(state => state.auth);

    // Find topic to edit
    const topic = topics.find(t => t.topicID === parseInt(topicID || '0'));

    const [title, setTitle] = useState(topic ? topic.title : '');
    const [description, setDescription] = useState(topic ? topic.description : '');

    useEffect(() => {
        if (topic) {
            setTitle(topic.title);
            setDescription(topic.description);
        }
    }, [topic]);

    // Check if user is author of topic
    const isAuthor = topic && topic.createdBy === userID;

    const handleTitleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setTitle(e.target.value);

        if (submitError) {
            dispatch(clearError());
        }
    };

    const handleDescriptionChange = (e: ChangeEvent<HTMLInputElement>) => {
        setDescription(e.target.value);

        if (submitError) {
            dispatch(clearError());
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        if (!topicID || !title.trim() || !description.trim()) return;

        const result = await dispatch(
            updateTopic({
                topicID: parseInt(topicID),
                title: title.trim(),
                description: description.trim()
            })
        );

        if (updateTopic.fulfilled.match(result)) {
            // Redirect to topic page on successful update
            navigate(`/topics/${topicID}`);
        }
    };

    if (loading) {
        return (
            <Container
                sx={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    minHeight: '60vh',
                }}
            >
                <CircularProgress />
            </Container>
        );
    }

    if (!topic) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">Topic not found.</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(`/topics`)}
                    sx={{ mt: 2 }}
                >
                    Back to Topics
                </Button>
            </Container>
        );
    }

    if (!isAuthor) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">You are not authorised to edit this topic.</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(`/topics/${topicID}`)}
                    sx={{ mt: 2 }}
                >
                    Back to Topic
                </Button>
            </Container>
        );
    }

    return (
        <Container
            maxWidth="md"
            sx={{
                mt: 4, 
                mb: 4,
            }}
        >
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate(`/topics/${topicID}`)}
                sx={{ mb: 3 }}
            >
                Back to Topic
            </Button>

            <Paper elevation={2} sx={{ p: 3 }}>
                <Typography
                    variant="h4"
                    component="h1"
                    gutterBottom
                >
                    Edit Topic
                </Typography>

                {submitError && (
                    <Alert severity="error" sx={{ mb: 2 }}>{submitError}</Alert>
                )}

                <form onSubmit={handleSubmit}>
                    <TextField 
                        fullWidth
                        label="Title"
                        value={title}
                        onChange={handleTitleChange}
                        disabled={submitting}
                        required
                        sx={{ mb: 2 }}
                    />

                    <TextField 
                        fullWidth
                        label="Description"
                        value={description}
                        onChange={handleDescriptionChange}
                        disabled={submitting}
                        required
                        multiline
                        minRows={3}
                        sx={{ mb: 2 }}
                    />

                    <Box
                        sx={{
                            display: 'flex',
                            gap: 2,
                        }}
                    >
                        <Button
                            type="submit"
                            variant="contained"
                            disabled={submitting || !title.trim() || !description.trim()}
                        >
                            {submitting ? <CircularProgress size={24} /> : 'Save Changes'}
                        </Button>    

                        <Button
                            variant="outlined"
                            onClick={() => navigate(`/topics/${topicID}`)}
                            disabled={submitting}
                        >
                            Cancel
                        </Button>
                    </Box>
                </form>
            </Paper>
        </Container>
    );
        
}