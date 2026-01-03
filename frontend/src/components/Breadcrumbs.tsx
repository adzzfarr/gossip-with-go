import { Home } from "@mui/icons-material";
import { Breadcrumbs, Link, Typography } from "@mui/material";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { useAppSelector } from "../hooks/redux";

export default function ForumBreadcrumbs() {
    const navigate = useNavigate();
    const location = useLocation();
    const { topicID, postID } = useParams<{ topicID?: string; postID?: string }>();

    // Get topics and current post from Redux state
    const { topics } = useAppSelector(state => state.topics);
    const { currentPost } = useAppSelector(state => state.posts);

    // Get current topic
    const currentTopic = topicID
        ? topics.find(topic => topic.topicID === parseInt(topicID))
        : currentPost
            ? topics.find(topic => topic.topicID === currentPost.topicID)
            : null;

    // Build breadcrumb items
    const getBreadcrumbItems = () => {
        const items: { label: string; path: string; icon?: React.ReactNode }[] = [
            {
                label: 'Topics',
                path: '/topics',
                icon: <Home sx={{ mr: 0.5 }} fontSize="inherit"/>
            }
        ]

        // On TopicPostsPage
        if (currentTopic) {
            items.push({
                label: currentTopic.title,
                path: `/topics/${currentTopic.topicID}`,
            });
        }

        // On CreatePostPage
        if (location.pathname.includes('create-post')) {
            items.push({
                label: 'Create Post',
                path: location.pathname,
            });
        }

        // On PostPage
        if (postID && currentPost) {
            items.push({
                label: currentPost.title,
                path: `/topics/${topicID}/posts/${postID}`,
            });
        }

        return items;
    }

    const breadcrumbItems = getBreadcrumbItems();

    if (location.pathname === '/topics') {
        // No breadcrumbs on TopicsPage (Home)
        return null; 
    }

    return (
        <Breadcrumbs
            aria-label="breadcrumb"
            sx={{ mb: 2 }}
        >
            {breadcrumbItems.map(
                (item, index) => {
                    const isLast = index === breadcrumbItems.length - 1;

                    return isLast
                        ? (
                            <Typography
                                key={item.path}
                                color="text.primary"
                                sx={{ 
                                    display: 'flex', 
                                    alignItems: 'center' 
                                }}
                            >
                                {item.icon}
                                {item.label}
                            </Typography>
                        )
                        : (
                            <Link
                                key={item.path}
                                underline="hover"
                                color="inherit"
                                onClick={() => navigate(item.path)}
                                sx={{ 
                                    cursor: 'pointer', 
                                    display: 'flex', 
                                    alignItems: 'center' 
                                }}    
                            >
                                {item.icon}
                                {item.label}
                            </Link>
                    );
                }
            )}
        </Breadcrumbs>
    );
}