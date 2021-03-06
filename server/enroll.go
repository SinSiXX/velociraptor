/*
   Velociraptor - Hunting Evil
   Copyright (C) 2019 Velocidex Innovations.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package server

import (
	"context"
	"fmt"

	"github.com/Velocidex/ordereddict"
	crypto_proto "www.velocidex.com/golang/velociraptor/crypto/proto"
	"www.velocidex.com/golang/velociraptor/logging"
	"www.velocidex.com/golang/velociraptor/result_sets"
	"www.velocidex.com/golang/velociraptor/services"
)

func enroll(
	ctx context.Context,
	server *Server,
	csr *crypto_proto.Certificate) error {

	if csr.GetType() != crypto_proto.Certificate_CSR || csr.Pem == nil {
		return nil
	}

	client_id, err := server.manager.AddCertificateRequest(csr.Pem)
	if err != nil {
		logger := logging.GetLogger(server.config, &logging.FrontendComponent)
		logger.Error(fmt.Sprintf("While enrolling %v: %v", client_id, err))
		return err
	}

	path_manager := result_sets.NewArtifactPathManager(
		server.config, "server" /* client_id */, "", "Server.Internal.Enrollment")

	return services.GetJournal().PushRows(path_manager,
		[]*ordereddict.Dict{ordereddict.NewDict().Set("ClientId", client_id)})
}
