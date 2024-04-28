package services

import (
	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"

	dtos "github.com/sm888sm/halten-backend/board-service/internal/models"
)

func convertMembersToProto(members []*dtos.BoardMemberDTO) []*pb_board.BoardMember {
	var protoMembers []*pb_board.BoardMember
	for _, member := range members {
		protoMember := &pb_board.BoardMember{
			UserID: member.ID,
			Role:   member.Role,
			// Assuming User model has Username and Fullname fields
			Username: member.Username,
			Fullname: member.Fullname,
		}
		protoMembers = append(protoMembers, protoMember)
	}
	return protoMembers
}

func convertListsToProto(lists []*pb_list.List) []*pb_board.List {
	var boardLists []*pb_board.List
	for _, list := range lists {
		boardList := &pb_board.List{
			ListID:   list.ListID,
			BoardID:  list.BoardID,
			Name:     list.Name,
			Position: list.Position,
		}
		boardLists = append(boardLists, boardList)
	}
	return boardLists
}

func convertCardsToProto(cards []*pb_card.CardMeta) []*pb_board.CardMeta {
	var boardLists []*pb_board.CardMeta

	for _, card := range cards {
		boardList := &pb_board.CardMeta{
			CardID:          card.CardID,
			BoardID:         card.BoardID,
			ListID:          card.ListID,
			Name:            card.Name,
			Position:        card.Position,
			Labels:          card.Labels,
			Members:         card.Members,
			TotalAttachment: card.TotalAttachment,
			TotalComment:    card.TotalComment,
		}

		boardLists = append(boardLists, boardList)
	}

	return boardLists
}

func convertLabelsToProto(labels []*dtos.LabelDTO) []*pb_board.Label {
	var labelsProto []*pb_board.Label

	for _, label := range labels {
		labelProto := &pb_board.Label{
			LabelID: label.ID,
			Name:    label.Name,
			Color:   label.Color,
			BoardID: label.BoardID,
		}

		labelsProto = append(labelsProto, labelProto)
	}

	return labelsProto
}
