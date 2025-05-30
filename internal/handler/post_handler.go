package handler

import (
	"fmt"
	"forum-project/internal/application"
	"forum-project/internal/auth"
	"forum-project/internal/model"
	"net/http"
	"strconv"
)

type PostHandler struct {
	app *application.App
}

func NewPostHandler(app *application.App) *PostHandler {
	return &PostHandler{app: app}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Responder.NotFound(rw, "Post Not Found", err)
		return
	}

	viewData := new(model.Page)
	viewData.IsAuthor = false

	claims, err := auth.GetClaimsFromRequest(r, p.app.Config.JWT.Secret)

	if err == nil {
		user := claims["user"].(map[string]interface{})
		userIDFloat := user["id"].(float64)
		userIDInt := int(userIDFloat)

		isAuthor := p.app.PostService.VerifyPostAuthor(post, userIDInt)
		if isAuthor {
			viewData.IsAuthor = true
		}
	}

	data := make(map[string]any)
	data["post"] = post

	viewData.Data = data

	err = p.app.Templates.Render(rw, r, "post.page", viewData)
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.app.TopicService.GetTopicByID(id)
	if err != nil {
		p.app.Responder.NotFound(rw, "Topic Not Found", err)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.app.Templates.Render(rw, r, "create-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostCreatePost(rw http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	topicID := r.PostFormValue("topic_id")
	topicIDInt, err := strconv.Atoi(topicID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Topic ID", err)
		return
	}

	user := r.Context().Value("user").(*model.AuthorizedUser)
	userID := user.ID
	userName := user.Username

	err = p.app.PostService.CreatePost(title, content, topicIDInt, userID, userName)
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to create post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%d", topicIDInt)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Responder.NotFound(rw, "Post Not Found", err)
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.app.Templates.Render(rw, r, "edit-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.app.Responder.NotFound(rw, "Topic Not Found", err)
		return
	}

	err = p.app.PostService.EditPost(title, content, id)
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to edit post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetDeletePost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Responder.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.app.Responder.NotFound(rw, "Topic Not Found", err)
		return
	}

	err = p.app.PostService.DeletePost(id)
	if err != nil {
		p.app.Responder.InternalServer(rw, "Unable to delete post", err)
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}
