import { Routes, Route, useNavigate, Navigate } from 'react-router-dom';
import { 
  Box, 
  Container, 
  Typography, 
  Button, 
  Card, 
  CardContent, 
  CardActions,
  Stack,
  Divider 
} from '@mui/material';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import { useAppDispatch, useAppSelector } from './hooks/redux';
import { logoutUser } from './features/auth/authSlice';
import TopicsPage from './pages/TopicsPage';
import TopicPostsPage from './pages/TopicPostsPage';
import PostPage from './pages/PostPage';
import CreatePostPage from './pages/CreatePostPage';
import Layout from './components/Layout';
import EditPostPage from './pages/EditPostPage';

function ThemePreview() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(logoutUser());
    navigate('/login');
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Stack spacing={4}>
        {/* Typography Section */}
        <Box>
          <Button variant='contained' onClick={handleLogout}>Logout</Button>
          <Typography variant="h1" gutterBottom>
            h1. Gossip with Go
          </Typography>
          <Typography variant="h2" gutterBottom>
            h2. Welcome to the Forum
          </Typography>
          <Typography variant="h3" gutterBottom>
            h3. Topic Discussions
          </Typography>
          <Typography variant="h4" gutterBottom>
            h4. Latest Posts
          </Typography>
          <Typography variant="h5" gutterBottom>
            h5. Community Comments
          </Typography>
          <Typography variant="h6" gutterBottom>
            h6. User Information
          </Typography>
          <Typography variant="body1" paragraph>
            Body 1: This is the default body text. Lorem ipsum dolor sit amet, 
            consectetur adipiscing elit. Material UI makes it easy to create 
            beautiful, responsive interfaces.
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Body 2: Secondary text with slightly smaller font size and muted color.
          </Typography>
        </Box>

        <Divider />

        {/* Buttons Section */}
        <Box>
          <Typography variant="h4" gutterBottom>
            MuiButton Examples
          </Typography>
          <Stack direction="row" spacing={2} flexWrap="wrap" sx={{ gap: 2 }}>
            <Button variant="contained">
              Contained Button
            </Button>
            <Button variant="contained" color="secondary">
              Secondary Button
            </Button>
            <Button variant="outlined">
              Outlined Button
            </Button>
            <Button variant="text">
              Text Button
            </Button>
            <Button variant="contained" disabled>
              Disabled Button
            </Button>
            <Button variant="contained" size="small">
              Small Button
            </Button>
            <Button variant="contained" size="large">
              Large Button
            </Button>
          </Stack>
        </Box>

        <Divider />

        {/* Cards Section */}
        <Box>
          <Typography variant="h4" gutterBottom>
            MuiCard Examples
          </Typography>
          <Stack spacing={3}>
            {/* Card 1 - Topic Card Example */}
            <Card>
              <CardContent>
                <Typography variant="h5" component="div" gutterBottom>
                  Example Topic: Technology Discussion
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  This is a sample topic card showcasing how topics might appear 
                  in the forum. Topics contain multiple posts and discussions.
                </Typography>
                <Typography variant="caption" display="block" sx={{ mt: 2 }}>
                  Created by: username • 2 hours ago
                </Typography>
              </CardContent>
              <CardActions>
                <Button size="small">View Posts</Button>
                <Button size="small">Edit</Button>
                <Button size="small" color="error">Delete</Button>
              </CardActions>
            </Card>

            {/* Card 2 - Post Card Example */}
            <Card elevation={3}>
              <CardContent>
                <Typography variant="h6" component="div" gutterBottom>
                  Example Post: Getting Started with Go
                </Typography>
                <Typography variant="body2" paragraph>
                  This is a sample post card with elevation. Posts belong to topics 
                  and can have multiple comments. The elevated card creates a 
                  nice shadow effect.
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Posted in: Technology Discussion • 15 comments
                </Typography>
              </CardContent>
              <CardActions>
                <Button size="small" variant="outlined">
                  Read More
                </Button>
                <Button size="small" variant="text">
                  Comment
                </Button>
              </CardActions>
            </Card>

            {/* Card 3 - Comment Card Example */}
            <Card variant="outlined">
              <CardContent>
                <Typography variant="body1" paragraph>
                  This is an example comment card with outlined variant. Comments 
                  are shorter and typically have less visual prominence than posts.
                </Typography>
                <Typography variant="caption" display="block">
                  By: commenter123 • 30 minutes ago
                </Typography>
              </CardContent>
              <CardActions>
                <Button size="small">Edit</Button>
                <Button size="small" color="error">Delete</Button>
              </CardActions>
            </Card>
          </Stack>
        </Box>

        <Divider />

        {/* Color Palette Section */}
        <Box>
          <Typography variant="h4" gutterBottom>
            Color Palette
          </Typography>
          <Stack direction="row" spacing={2} flexWrap="wrap" sx={{ gap: 2 }}>
            <Box
              sx={{
                width: 120,
                height: 120,
                backgroundColor: 'primary.main',
                borderRadius: 2,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography color="white" fontWeight="bold">
                Primary
              </Typography>
            </Box>
            <Box
              sx={{
                width: 120,
                height: 120,
                backgroundColor: 'secondary.main',
                borderRadius: 2,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography color="white" fontWeight="bold">
                Secondary
              </Typography>
            </Box>
            <Box
              sx={{
                width: 120,
                height: 120,
                backgroundColor: 'background.paper',
                borderRadius: 2,
                border: '1px solid',
                borderColor: 'divider',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography fontWeight="bold">
                Paper
              </Typography>
            </Box>
            <Box
              sx={{
                width: 120,
                height: 120,
                backgroundColor: 'error.main',
                borderRadius: 2,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography color="white" fontWeight="bold">
                Error
              </Typography>
            </Box>
          </Stack>
        </Box>
      </Stack>
    </Container>
  );
}

function App() {
  const { isAuthenticated } = useAppSelector(state => state.auth);

  return (
    <Routes>
      {/* Public Routes (No Layout) */}
      <Route path='/login' element={<LoginPage />} />
      <Route path='/register' element={<RegisterPage />} />
      <Route path='/theme' element={<ThemePreview />} />

      {/* Protected Routes (With Layout) */}
      <Route 
        path='/topics' 
        element={
          isAuthenticated 
            ? (<Layout>
                <TopicsPage />
              </Layout>) 
            : (<Navigate to="/login" replace />)
        } 
      />

      <Route 
        path='/topics/:topicID' 
        element={
          isAuthenticated 
            ? (<Layout>
                <TopicPostsPage />
              </Layout>) 
            : (<Navigate to="/login" replace />)
        } 
      />

      <Route 
        path='/topics/:topicID/create-post'
        element={
          isAuthenticated 
            ? (<Layout>
                <CreatePostPage />
              </Layout>) 
            : (<Navigate to="/login" replace />)
        }
      />

      <Route
        path='/posts/:postID'
        element={
          isAuthenticated 
            ? (<Layout>
                <PostPage />
              </Layout>) 
            : (<Navigate to="/login" replace />)
        }
      />

      <Route
        path='/posts/:postID/edit'
        element={
          isAuthenticated 
            ? (<Layout>
                <EditPostPage />
              </Layout>) 
            : (<Navigate to="/login" replace />)
        }
      />

      {/* Default Route */}
      <Route 
        path='/' 
        element={isAuthenticated ? <Navigate to="/topics" replace /> : <Navigate to="/login" replace />} 
      />
    </Routes>
  );
}

export default App;